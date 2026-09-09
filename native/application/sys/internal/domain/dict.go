package sys

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/lightning/application/sys/internal/domain/model"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DictService 字典服务接口
type DictService interface {
	// Create 创建字典
	Create(ctx context.Context, req *CreateDictRequest) error

	// Update 更新字典
	Update(ctx context.Context, req *UpdateDictRequest) error

	// Delete 删除字典
	Delete(ctx context.Context, id int64) error

	// BatchDelete 批量删除字典
	BatchDelete(ctx context.Context, ids []int64) error

	// GetById 根据ID查询字典
	GetById(ctx context.Context, id int64) (*model.DictData, error)

	// Page 分页查询字典列表
	Page(ctx context.Context, req *PageDictRequest) (*pagination.Page[model.DictData], error)

	// GetByType 根据字典类型获取字典列表
	GetByType(ctx context.Context, dictType string) ([]model.DictData, error)

	// GetByTypeAndParent 根据字典类型和父ID获取子字典列表
	GetByTypeAndParent(ctx context.Context, dictType string, parentId int64) ([]model.DictData, error)

	// GetDictLabel 根据字典类型和键值获取标签
	GetDictLabel(ctx context.Context, dictType, dictValue string) (string, error)

	// GetDictValue 根据字典类型和标签获取键值
	GetDictValue(ctx context.Context, dictType, dictLabel string) (string, error)
}

type dictService struct {
	db     *gorm.DB
	logger logging.Logger
}

// NewDictService 创建字典服务实例
func NewDictService(db *gorm.DB, logger logging.Logger) DictService {
	return &dictService{
		db:     db,
		logger: logger,
	}
}

// Create 创建字典
func (s *dictService) Create(ctx context.Context, req *CreateDictRequest) error {
	// 检查字典值是否已存在（同类型下）
	exists, err := (&model.DictData{}).CheckDictValueExists(s.db, req.DictType, req.DictValue)
	if err != nil {
		s.logger.Error("检查字典值失败", zap.Error(err))
		return fmt.Errorf("检查字典值失败: %w", err)
	}
	if exists {
		return fmt.Errorf("字典值已存在: %s", req.DictValue)
	}

	// 如果指定了父字典，检查父字典是否存在
	if req.ParentId > 0 {
		parent, err := (&model.DictData{}).FindByID(s.db, req.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("父字典不存在")
			}
			s.logger.Error("查询父字典失败", zap.Error(err))
			return fmt.Errorf("查询父字典失败: %w", err)
		}

		// 父字典和子字典必须是同一类型
		if parent.DictType != req.DictType {
			return fmt.Errorf("父字典类型不匹配")
		}
	}

	// 创建字典
	dict := &model.DictData{
		ParentId:  req.ParentId,
		DictType:  req.DictType,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		Sort:      req.Sort,
		IsDefault: req.IsDefault,
		Status:    req.Status,
		Remark:    req.Remark,
		CreateBy:  req.CreateBy,
		UpdateBy:  req.UpdateBy,
	}

	if err := dict.Create(s.db); err != nil {
		s.logger.Error("创建字典失败", zap.Error(err))
		return fmt.Errorf("创建字典失败: %w", err)
	}

	s.logger.Info("创建字典成功",
		zap.Int64("id", dict.ID),
		zap.String("dictType", dict.DictType),
		zap.String("dictLabel", dict.DictLabel))

	return nil
}

// Update 更新字典
func (s *dictService) Update(ctx context.Context, req *UpdateDictRequest) error {
	// 检查字典是否存在
	existingDict, err := (&model.DictData{}).FindByID(s.db, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("字典不存在")
		}
		s.logger.Error("查询字典失败", zap.Error(err))
		return fmt.Errorf("查询字典失败: %w", err)
	}

	// 检查字典值是否被其他字典占用
	if req.DictValue != existingDict.DictValue {
		exists, err := (&model.DictData{}).CheckDictValueExistsExcludingSelf(
			s.db, req.ID, req.DictType, req.DictValue,
		)
		if err != nil {
			s.logger.Error("检查字典值失败", zap.Error(err))
			return fmt.Errorf("检查字典值失败: %w", err)
		}
		if exists {
			return fmt.Errorf("字典值已被占用: %s", req.DictValue)
		}
	}

	// 不能将自己设置为父字典
	if req.ParentId == req.ID {
		return fmt.Errorf("不能将自己设置为父字典")
	}

	// 如果修改了父字典，检查父字典是否存在
	if req.ParentId > 0 && req.ParentId != existingDict.ParentId {
		parent, err := (&model.DictData{}).FindByID(s.db, req.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("父字典不存在")
			}
			s.logger.Error("查询父字典失败", zap.Error(err))
			return fmt.Errorf("查询父字典失败: %w", err)
		}

		// 父字典和子字典必须是同一类型
		if parent.DictType != req.DictType {
			return fmt.Errorf("父字典类型不匹配")
		}
	}

	// 更新字典
	updates := map[string]interface{}{
		"parent_id":  req.ParentId,
		"dict_type":  req.DictType,
		"dict_label": req.DictLabel,
		"dict_value": req.DictValue,
		"sort":       req.Sort,
		"is_default": req.IsDefault,
		"status":     req.Status,
		"remark":     req.Remark,
		"update_by":  req.UpdateBy,
	}

	if err := existingDict.Update(s.db, req.ID, updates); err != nil {
		s.logger.Error("更新字典失败", zap.Error(err))
		return fmt.Errorf("更新字典失败: %w", err)
	}

	s.logger.Info("更新字典成功", zap.Int64("id", req.ID))
	return nil
}

// Delete 删除字典（级联删除子字典）
func (s *dictService) Delete(ctx context.Context, id int64) error {
	// 检查字典是否存在
	_, err := (&model.DictData{}).FindByID(s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("字典不存在")
		}
		s.logger.Error("查询字典失败", zap.Error(err))
		return fmt.Errorf("查询字典失败: %w", err)
	}

	// 使用递归CTE查询所有需要删除的ID并批量删除
	sql := `
		WITH RECURSIVE dict_tree AS (
			SELECT id FROM s_dict_data WHERE id = ?
			UNION ALL
			SELECT d.id FROM s_dict_data d
			INNER JOIN dict_tree dt ON d.parent_id = dt.id
		)
		DELETE FROM s_dict_data WHERE id IN (SELECT id FROM dict_tree)
	`

	if err := s.db.Exec(sql, id).Error; err != nil {
		s.logger.Error("删除字典失败", zap.Error(err))
		return fmt.Errorf("删除字典失败: %w", err)
	}

	s.logger.Info("删除字典成功", zap.Int64("id", id))
	return nil
}

// BatchDelete 批量删除字典（级联删除子字典）
func (s *dictService) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("字典ID列表不能为空")
	}

	// 使用递归CTE查询所有需要删除的ID并批量删除
	sql := `
		WITH RECURSIVE dict_tree AS (
			SELECT id FROM s_dict_data WHERE id = ANY(?)
			UNION ALL
			SELECT d.id FROM s_dict_data d
			INNER JOIN dict_tree dt ON d.parent_id = dt.id
		)
		DELETE FROM s_dict_data WHERE id IN (SELECT id FROM dict_tree)
	`

	if err := s.db.Exec(sql, ids).Error; err != nil {
		s.logger.Error("批量删除字典失败", zap.Error(err))
		return fmt.Errorf("批量删除字典失败: %w", err)
	}

	s.logger.Info("批量删除字典成功", zap.Int("count", len(ids)))
	return nil
}

// GetById 根据ID查询字典
func (s *dictService) GetById(ctx context.Context, id int64) (*model.DictData, error) {
	dict, err := (&model.DictData{}).FindByID(s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("字典不存在")
		}
		s.logger.Error("查询字典失败", zap.Error(err))
		return nil, fmt.Errorf("查询字典失败: %w", err)
	}
	return dict, nil
}

// Page 分页查询字典列表
func (s *dictService) Page(ctx context.Context, req *PageDictRequest) (*pagination.Page[model.DictData], error) {
	// 构建查询条件
	query := s.db.Model(&model.DictData{})

	// 仅查询顶级字典（parent_id = 0 或 NULL）
	query = query.Where("parent_id = 0 OR parent_id IS NULL")

	// 添加条件过滤
	if req.DictType != "" {
		query = query.Where("dict_type = ?", req.DictType)
	}
	if req.DictLabel != "" {
		query = query.Where("dict_label LIKE ?", "%"+req.DictLabel+"%")
	}
	if req.Status >= 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 添加默认排序
	if req.PageQuery.OrderByColumn == "" {
		query = query.Order("sort ASC, id DESC")
	}

	// 使用 Paginator 执行分页查询
	page, err := pagination.New[model.DictData](query, &req.PageQuery).Find()
	if err != nil {
		s.logger.Error("查询字典列表失败",
			zap.Error(err),
			zap.String("dictType", req.DictType),
			zap.String("dictLabel", req.DictLabel),
			zap.Int32("status", req.Status))
		return nil, fmt.Errorf("查询字典列表失败: %w", err)
	}

	return page, nil
}

// GetByType 根据字典类型获取字典列表
func (s *dictService) GetByType(ctx context.Context, dictType string) ([]model.DictData, error) {
	dicts, err := (&model.DictData{}).FindByType(s.db, dictType)
	if err != nil {
		s.logger.Error("查询字典列表失败", zap.Error(err))
		return nil, fmt.Errorf("查询字典列表失败: %w", err)
	}
	return dicts, nil
}

// GetByTypeAndParent 根据字典类型和父ID获取子字典列表
func (s *dictService) GetByTypeAndParent(ctx context.Context, dictType string, parentId int64) ([]model.DictData, error) {
	dicts, err := (&model.DictData{}).FindByTypeAndParent(s.db, dictType, parentId)
	if err != nil {
		s.logger.Error("查询子字典列表失败", zap.Error(err))
		return nil, fmt.Errorf("查询子字典列表失败: %w", err)
	}
	return dicts, nil
}

// GetDictLabel 根据字典类型和键值获取标签
func (s *dictService) GetDictLabel(ctx context.Context, dictType, dictValue string) (string, error) {
	label, err := (&model.DictData{}).GetDictLabel(s.db, dictType, dictValue)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("字典不存在")
		}
		return "", fmt.Errorf("查询字典标签失败: %w", err)
	}
	return label, nil
}

// GetDictValue 根据字典类型和标签获取键值
func (s *dictService) GetDictValue(ctx context.Context, dictType, dictLabel string) (string, error) {
	value, err := (&model.DictData{}).GetDictValue(s.db, dictType, dictLabel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("字典不存在")
		}
		return "", fmt.Errorf("查询字典键值失败: %w", err)
	}
	return value, nil
}
