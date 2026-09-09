package data

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	entpkg "github.com/gcc798/lightning/kratos/application/sys-rpc/ent"
)

func (r *Resources) attachmentEntityToItem(ctx context.Context, item *entpkg.Attachment) *v1.AttachmentItem {
	if item == nil {
		return nil
	}
	accessURL := strings.TrimSpace(item.AccessURL)
	if item.IsPublic && r != nil && r.Storage != nil {
		if resolved, err := r.Storage.AttachmentURL(ctx, item, 0); err == nil && strings.TrimSpace(resolved) != "" {
			accessURL = resolved
		}
	}
	return &v1.AttachmentItem{
		AttachmentId:  item.ID,
		EnvId:         item.EnvID,
		FileName:      item.FileName,
		FileKey:       item.FileKey,
		FileSize:      item.FileSize,
		FileType:      item.FileType,
		FileExt:       item.FileExt,
		BusinessType:  item.BusinessType,
		BusinessId:    item.BusinessID,
		BusinessField: item.BusinessField,
		IsPublic:      item.IsPublic,
		AccessUrl:     accessURL,
		Metadata:      metadataString(item.Metadata),
		Status:        item.Status,
		ExpireTime:    formatOptionalTime(item.ExpireTime),
		CreateBy:      item.CreateBy,
		CreateTime:    formatTime(item.CreateTime),
		UpdateTime:    formatTime(item.UpdateTime),
	}
}

func (r *Resources) SaveAttachment(ctx context.Context, req *v1.UploadFileRequest) (*v1.AttachmentItem, error) {
	if strings.TrimSpace(req.GetFileName()) == "" {
		return nil, errors.New("file name is required")
	}
	if len(req.GetContent()) == 0 {
		return nil, errors.New("file content is empty")
	}
	if r.Storage == nil {
		return nil, errors.New("storage manager is not initialized")
	}
	fileName := filepath.Base(strings.TrimSpace(req.GetFileName()))
	if fileName == "." || fileName == string(filepath.Separator) {
		return nil, errors.New("file name is invalid")
	}
	now := time.Now()
	attachmentID := nextID()
	fileType := strings.TrimSpace(req.GetFileType())
	if fileType == "" {
		fileType = http.DetectContentType(req.GetContent())
	}
	fileKey := filepath.ToSlash(filepath.Join(now.Format("20060102"), fmt.Sprintf("%d_%s", attachmentID, fileName)))
	env, err := r.Storage.Upload(ctx, req.GetEnvCode(), fileKey, bytes.NewReader(req.GetContent()), int64(len(req.GetContent())), fileType)
	if err != nil {
		return nil, err
	}
	createBy := currentOperatorID(ctx)
	item, err := r.Ent.Attachment.Create().
		SetID(attachmentID).
		SetEnvID(env.ID).
		SetFileName(fileName).
		SetFileKey(fileKey).
		SetFileSize(int64(len(req.GetContent()))).
		SetFileType(fileType).
		SetFileExt(strings.TrimPrefix(filepath.Ext(fileName), ".")).
		SetBusinessType("").
		SetBusinessID("").
		SetBusinessField("").
		SetIsPublic(false).
		SetAccessURL("").
		SetMetadata(map[string]any{}).
		SetStatus(0).
		SetCreateBy(createBy).
		SetCreateTime(now).
		SetUpdateTime(now).
		Save(ctx)
	if err != nil {
		_ = r.Storage.Delete(ctx, env.ID, fileKey)
		return nil, err
	}
	return r.attachmentEntityToItem(ctx, item), nil
}

func (r *Resources) GetAttachment(ctx context.Context, id int64) (*v1.AttachmentItem, error) {
	item, err := r.loadAttachmentEntity(ctx, id)
	if err != nil || item == nil {
		return nil, err
	}
	return r.attachmentEntityToItem(ctx, item), nil
}

func (r *Resources) BindAttachment(ctx context.Context, req *v1.BindAttachmentRequest) error {
	attachmentID := req.GetAttachmentId()
	item, err := r.loadAttachmentEntity(ctx, attachmentID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("attachment not found")
	}
	businessType := strings.TrimSpace(req.GetBusinessType())
	businessID := strings.TrimSpace(req.GetBusinessId())
	if businessType == "" || businessID == "" {
		return errors.New("business type and business id are required")
	}
	metadata := item.Metadata
	if req.GetMetadata() != nil {
		metadata = protoStructToMap(req.GetMetadata())
	}
	candidate := *item
	candidate.BusinessType = businessType
	candidate.BusinessID = businessID
	candidate.BusinessField = req.GetBusinessField()
	candidate.IsPublic = req.GetIsPublic()
	candidate.Metadata = metadata
	accessURL := ""
	if req.GetIsPublic() {
		if r.Storage == nil {
			return errors.New("storage manager is not initialized")
		}
		resolvedURL, err := r.Storage.AttachmentURL(ctx, &candidate, 0)
		if err != nil {
			return err
		}
		accessURL = resolvedURL
	}
	expireTime, err := parseAttachmentExpireTime(req.GetExpireTime())
	if err != nil {
		return err
	}
	updater := r.Ent.Attachment.UpdateOneID(attachmentID).
		SetBusinessType(businessType).
		SetBusinessID(businessID).
		SetBusinessField(req.GetBusinessField()).
		SetIsPublic(req.GetIsPublic()).
		SetAccessURL(accessURL).
		SetMetadata(metadata).
		SetUpdateTime(time.Now())
	if expireTime != nil {
		updater.SetExpireTime(*expireTime)
	} else {
		updater.ClearExpireTime()
	}
	_, err = updater.Save(ctx)
	return err
}

func (r *Resources) ListAttachments(ctx context.Context, businessType, businessID string) ([]*v1.AttachmentItem, error) {
	items, err := r.Ent.Attachment.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.AttachmentItem, 0)
	for _, item := range items {
		if item.DeletedAt != nil {
			continue
		}
		if businessType != "" && item.BusinessType != businessType {
			continue
		}
		if businessID != "" && item.BusinessID != businessID {
			continue
		}
		out = append(out, r.attachmentEntityToItem(ctx, item))
	}
	return out, nil
}

func (r *Resources) PageAttachments(ctx context.Context, req *v1.PageAttachmentsRequest) (*v1.PageAttachmentsReply, error) {
	items, err := r.Ent.Attachment.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.AttachmentItem, 0)
	for _, item := range items {
		if item.DeletedAt != nil {
			continue
		}
		if req.GetFileName() != "" && !strings.Contains(item.FileName, req.GetFileName()) {
			continue
		}
		if req.GetFileType() != "" && item.FileType != req.GetFileType() {
			continue
		}
		if req.GetBusinessType() != "" && item.BusinessType != req.GetBusinessType() {
			continue
		}
		filtered = append(filtered, r.attachmentEntityToItem(ctx, item))
	}
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageAttachmentsReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}

func (r *Resources) DownloadAttachment(ctx context.Context, id int64) (*v1.AttachmentDownloadReply, error) {
	item, err := r.loadAttachmentEntity(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	if r.Storage == nil {
		return nil, errors.New("storage manager is not initialized")
	}
	reader, err := r.Storage.Download(ctx, item.EnvID, item.FileKey)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return &v1.AttachmentDownloadReply{FileName: item.FileName, Content: content}, nil
}

func (r *Resources) GetAttachmentURL(ctx context.Context, id int64, expiresSeconds int64) (*v1.AttachmentURLReply, error) {
	item, err := r.loadAttachmentEntity(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	if r.Storage == nil {
		return nil, errors.New("storage manager is not initialized")
	}
	expires := normalizeAttachmentURLExpires(expiresSeconds)
	url, err := r.Storage.AttachmentURL(ctx, item, time.Duration(expires)*time.Second)
	if err != nil {
		return nil, err
	}
	return &v1.AttachmentURLReply{AttachmentId: id, Url: url, Expires: expires}, nil
}

func (r *Resources) DeleteAttachments(ctx context.Context, ids ...int64) error {
	now := time.Now()
	for _, id := range ids {
		item, err := r.loadAttachmentEntity(ctx, id)
		if err != nil {
			return err
		}
		if item == nil {
			return errors.New("attachment not found")
		}
		if r.Storage != nil {
			if err := r.Storage.Delete(ctx, item.EnvID, item.FileKey); err != nil {
				return err
			}
		}
		if _, err := r.Ent.Attachment.UpdateOneID(id).SetDeletedAt(now).SetUpdateTime(now).SetStatus(1).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (r *Resources) loadAttachmentEntity(ctx context.Context, id int64) (*entpkg.Attachment, error) {
	item, err := r.Ent.Attachment.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return item, nil
}

func metadataString(metadata map[string]any) string {
	if len(metadata) == 0 {
		return ""
	}
	if value, ok := metadata["raw"].(string); ok && strings.TrimSpace(value) != "" {
		return value
	}
	return stringifyJSON(metadata)
}

func formatOptionalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return formatTime(*t)
}

func parseAttachmentExpireTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New("invalid expireTime")
}

func normalizeAttachmentURLExpires(value int64) int64 {
	if value <= 0 {
		return 3600
	}
	return value
}
