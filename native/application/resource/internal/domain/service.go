package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/platform/storage"
	"github.com/gcc798/lightning/internal/utils"
	"github.com/gcc798/lightning/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ErrAttachmentNotFound 附件不存在。
var ErrAttachmentNotFound = errors.New("附件不存在")

// BindInput 绑定附件到业务的入参。
type BindInput struct {
	BusinessType  string
	BusinessId    string
	BusinessField string
	IsPublic      bool
	Metadata      map[string]any
	ExpireTime    *utils.LocalTime
}

// PageQuery 附件分页查询条件。
type PageQuery struct {
	PageNum      int
	PageSize     int
	FileName     string
	FileType     string
	BusinessType string
}

// AttachmentService 附件领域服务。
type AttachmentService interface {
	UploadFile(ctx context.Context, r io.Reader, filename, contentType string, size int64) (*Attachment, error)
	BindToBusiness(ctx context.Context, attachmentId int64, in BindInput) error
	Download(ctx context.Context, attachmentId int64) (io.ReadCloser, string, error)
	Delete(ctx context.Context, attachmentId int64) error
	GetURL(ctx context.Context, attachmentId int64, expires time.Duration) (string, error)
	GetById(ctx context.Context, attachmentId int64) (*Attachment, error)
	ListByBusiness(ctx context.Context, businessType, businessId string) ([]*Attachment, error)
	Page(ctx context.Context, q PageQuery) (*pagination.Page[Attachment], error)
	CleanExpired(ctx context.Context) (cleaned, failed int64, err error)
}

type attachmentService struct {
	db      *gorm.DB
	storage storage.Storage
	logger  logger.Logger
}

// NewAttachmentService 创建附件服务实例。
func NewAttachmentService(db *gorm.DB, store storage.Storage, log logger.Logger) AttachmentService {
	return &attachmentService{db: db, storage: store, logger: log}
}

// UploadFile 上传文件（步骤1：只落存储与附件记录，尚未关联业务）。
func (s *attachmentService) UploadFile(ctx context.Context, r io.Reader, filename, contentType string, size int64) (*Attachment, error) {
	fileKey := generateFileKey(filename, "temp")

	if err := s.storage.Upload(ctx, fileKey, r, size); err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	attachment := &Attachment{
		FileName: filename,
		FileKey:  fileKey,
		FileSize: size,
		FileType: contentType,
		FileExt:  strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), ".")),
		Status:   StatusNormal,
	}

	if err := s.db.WithContext(ctx).Create(attachment).Error; err != nil {
		// 记录落库失败，回滚已上传的对象
		if delErr := s.storage.Delete(ctx, fileKey); delErr != nil {
			s.logger.Error("回滚存储对象失败", zap.String("fileKey", fileKey), zap.Error(delErr))
		}
		return nil, fmt.Errorf("保存附件记录失败: %w", err)
	}

	s.logger.Info("上传文件成功",
		zap.Int64("attachmentId", attachment.ID),
		zap.String("fileName", attachment.FileName),
		zap.Int64("fileSize", attachment.FileSize))

	return attachment, nil
}

// BindToBusiness 绑定附件到业务（步骤2）。
func (s *attachmentService) BindToBusiness(ctx context.Context, attachmentId int64, in BindInput) error {
	attachment, err := s.GetById(ctx, attachmentId)
	if err != nil {
		return err
	}

	// 公开文件生成永久访问 URL，失败不阻断绑定
	var accessUrl string
	if in.IsPublic {
		accessUrl, err = s.storage.GetURL(ctx, attachment.FileKey, 0)
		if err != nil {
			s.logger.Warn("获取访问 URL 失败", zap.Int64("attachmentId", attachmentId), zap.Error(err))
		}
	}

	var metadata *json.RawMessage
	if len(in.Metadata) > 0 {
		raw, err := json.Marshal(in.Metadata)
		if err != nil {
			return fmt.Errorf("元数据序列化失败: %w", err)
		}
		msg := json.RawMessage(raw)
		metadata = &msg
	}

	updates := map[string]any{
		"business_type":  in.BusinessType,
		"business_id":    in.BusinessId,
		"business_field": in.BusinessField,
		"is_public":      in.IsPublic,
		"access_url":     accessUrl,
		"metadata":       metadata,
	}
	if in.ExpireTime != nil {
		updates["expire_time"] = *in.ExpireTime
	}

	if err := s.db.WithContext(ctx).Model(&Attachment{}).
		Where("id = ?", attachmentId).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("绑定附件到业务失败: %w", err)
	}

	s.logger.Info("绑定附件到业务成功",
		zap.Int64("attachmentId", attachmentId),
		zap.String("businessType", in.BusinessType),
		zap.String("businessId", in.BusinessId))

	return nil
}

// Download 下载附件，返回内容流与原始文件名。
func (s *attachmentService) Download(ctx context.Context, attachmentId int64) (io.ReadCloser, string, error) {
	attachment, err := s.GetById(ctx, attachmentId)
	if err != nil {
		return nil, "", err
	}

	reader, err := s.storage.Download(ctx, attachment.FileKey)
	if err != nil {
		return nil, "", fmt.Errorf("下载文件失败: %w", err)
	}

	return reader, attachment.FileName, nil
}

// Delete 删除存储对象并软删除附件记录。
func (s *attachmentService) Delete(ctx context.Context, attachmentId int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var attachment Attachment
		if err := tx.Where("id = ?", attachmentId).First(&attachment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAttachmentNotFound
			}
			return fmt.Errorf("查询附件失败: %w", err)
		}

		// 存储对象删除失败不阻断记录标记，避免记录与对象都无法回收
		if err := s.storage.Delete(ctx, attachment.FileKey); err != nil {
			s.logger.Warn("删除存储对象失败", zap.String("fileKey", attachment.FileKey), zap.Error(err))
		}

		if err := tx.Model(&Attachment{}).
			Where("id = ?", attachmentId).
			Update("status", StatusDeleted).Error; err != nil {
			return fmt.Errorf("删除附件记录失败: %w", err)
		}

		s.logger.Info("删除附件成功", zap.Int64("attachmentId", attachmentId))
		return nil
	})
}

// GetURL 获取附件访问 URL，expires 为 0 表示永久。
func (s *attachmentService) GetURL(ctx context.Context, attachmentId int64, expires time.Duration) (string, error) {
	attachment, err := s.GetById(ctx, attachmentId)
	if err != nil {
		return "", err
	}

	if attachment.IsPublic && attachment.AccessUrl != "" && expires == 0 {
		return attachment.AccessUrl, nil
	}

	url, err := s.storage.GetURL(ctx, attachment.FileKey, expires)
	if err != nil {
		return "", fmt.Errorf("生成访问 URL 失败: %w", err)
	}

	return url, nil
}

// GetById 根据 ID 查询未删除的附件。
func (s *attachmentService) GetById(ctx context.Context, attachmentId int64) (*Attachment, error) {
	var attachment Attachment
	if err := s.db.WithContext(ctx).
		Where("id = ? AND status = ?", attachmentId, StatusNormal).
		First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttachmentNotFound
		}
		return nil, fmt.Errorf("查询附件失败: %w", err)
	}
	return &attachment, nil
}

// ListByBusiness 根据业务查询附件列表。
func (s *attachmentService) ListByBusiness(ctx context.Context, businessType, businessId string) ([]*Attachment, error) {
	var attachments []*Attachment
	if err := s.db.WithContext(ctx).
		Where("business_type = ? AND business_id = ? AND status = ?", businessType, businessId, StatusNormal).
		Order("create_time DESC").
		Find(&attachments).Error; err != nil {
		return nil, fmt.Errorf("查询附件列表失败: %w", err)
	}
	return attachments, nil
}

// Page 分页查询附件列表。
func (s *attachmentService) Page(ctx context.Context, q PageQuery) (*pagination.Page[Attachment], error) {
	query := s.db.WithContext(ctx).Model(&Attachment{}).Where("status = ?", StatusNormal)
	if q.FileName != "" {
		query = query.Where("file_name LIKE ?", "%"+q.FileName+"%")
	}
	if q.FileType != "" {
		query = query.Where("file_type = ?", q.FileType)
	}
	if q.BusinessType != "" {
		query = query.Where("business_type = ?", q.BusinessType)
	}

	page, err := pagination.New[Attachment](query, &pagination.PageQuery{
		PageNum:  q.PageNum,
		PageSize: q.PageSize,
	}).Find()
	if err != nil {
		return nil, fmt.Errorf("分页查询附件列表失败: %w", err)
	}
	return page, nil
}

// CleanExpired 清理已过期的附件。
func (s *attachmentService) CleanExpired(ctx context.Context) (int64, int64, error) {
	var attachments []*Attachment
	if err := s.db.WithContext(ctx).
		Where("expire_time IS NOT NULL AND expire_time < ? AND status = ?", time.Now(), StatusNormal).
		Find(&attachments).Error; err != nil {
		return 0, 0, fmt.Errorf("查询过期附件失败: %w", err)
	}

	var failed int64
	for _, attachment := range attachments {
		if err := s.Delete(ctx, attachment.ID); err != nil {
			failed++
			s.logger.Error("删除过期附件失败", zap.Int64("attachmentId", attachment.ID), zap.Error(err))
		}
	}

	s.logger.Info("清理过期附件完成", zap.Int("total", len(attachments)), zap.Int64("failed", failed))
	return int64(len(attachments)) - failed, failed, nil
}

// generateFileKey 生成存储 Key：{businessType}/{date}/{timestamp}_{filename}。
func generateFileKey(filename, businessType string) string {
	now := time.Now()
	return fmt.Sprintf("%s/%s/%d_%s", businessType, now.Format("20060102"), now.UnixNano(), filepath.Base(filename))
}
