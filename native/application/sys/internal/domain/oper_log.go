package sys

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/lightning/application/sys/internal/domain/model"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/utils"
	"github.com/gcc798/lightning/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OperLogService 操作日志服务接口
type OperLogService interface {
	// Create 创建操作日志
	Create(ctx context.Context, req *CreateOperLogRequest) error

	// Update 更新操作日志
	Update(ctx context.Context, req *UpdateOperLogRequest) error

	// Delete 删除操作日志
	Delete(ctx context.Context, id int64) error

	// BatchDelete 批量删除操作日志
	BatchDelete(ctx context.Context, ids []int64) error

	// GetById 根据ID查询操作日志
	GetById(ctx context.Context, id int64) (*model.OperLog, error)

	// Page 分页查询操作日志列表
	Page(ctx context.Context, req *PageOperLogRequest) (*pagination.Page[model.OperLog], error)

	// CleanOldLogs 清理指定天数之前的日志
	CleanOldLogs(ctx context.Context, days int) (int64, error)
}

type operLogService struct {
	db     *gorm.DB
	logger logging.Logger
}

// NewOperLogService 创建操作日志服务实例
func NewOperLogService(db *gorm.DB, logger logging.Logger) OperLogService {
	return &operLogService{
		db:     db,
		logger: logger,
	}
}

// Create 创建操作日志
func (s *operLogService) Create(ctx context.Context, req *CreateOperLogRequest) error {
	log := &model.OperLog{
		Title:         req.Title,
		BusinessType:  req.BusinessType,
		Method:        req.Method,
		RequestMethod: req.RequestMethod,
		DeviceType:    req.DeviceType,
		OperName:      req.OperName,
		OperUrl:       req.OperUrl,
		OperIp:        req.OperIp,
		OperLocation:  req.OperLocation,
		OperParam:     req.OperParam,
		JsonResult:    req.JsonResult,
		Status:        req.Status,
		ErrorMsg:      req.ErrorMsg,
		OperTime:      utils.Now(),
		CostTime:      req.CostTime,
		UserAgent:     req.UserAgent,
	}

	if err := log.Create(s.db); err != nil {
		s.logger.Error("创建操作日志失败", zap.Error(err))
		return fmt.Errorf("创建操作日志失败: %w", err)
	}

	s.logger.Info("创建操作日志成功",
		zap.Int64("id", log.ID),
		zap.String("title", log.Title),
		zap.String("operName", log.OperName))

	return nil
}

// Update 更新操作日志
func (s *operLogService) Update(ctx context.Context, req *UpdateOperLogRequest) error {
	// 检查日志是否存在
	_, err := (&model.OperLog{}).FindByID(s.db, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("操作日志不存在")
		}
		s.logger.Error("查询操作日志失败", zap.Error(err))
		return fmt.Errorf("查询操作日志失败: %w", err)
	}

	// 更新日志
	updates := map[string]interface{}{
		"title":          req.Title,
		"business_type":  req.BusinessType,
		"method":         req.Method,
		"request_method": req.RequestMethod,
		"device_type":    req.DeviceType,
		"oper_name":      req.OperName,
		"oper_url":       req.OperUrl,
		"oper_ip":        req.OperIp,
		"oper_location":  req.OperLocation,
		"oper_param":     req.OperParam,
		"json_result":    req.JsonResult,
		"status":         req.Status,
		"error_msg":      req.ErrorMsg,
		"cost_time":      req.CostTime,
		"user_agent":     req.UserAgent,
	}

	log := &model.OperLog{}
	if err := log.Update(s.db, req.ID, updates); err != nil {
		s.logger.Error("更新操作日志失败", zap.Error(err))
		return fmt.Errorf("更新操作日志失败: %w", err)
	}

	s.logger.Info("更新操作日志成功", zap.Int64("id", req.ID))
	return nil
}

// Delete 删除操作日志
func (s *operLogService) Delete(ctx context.Context, id int64) error {
	// 检查日志是否存在
	_, err := (&model.OperLog{}).FindByID(s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("操作日志不存在")
		}
		s.logger.Error("查询操作日志失败", zap.Error(err))
		return fmt.Errorf("查询操作日志失败: %w", err)
	}

	// 删除日志
	if err := (&model.OperLog{}).Delete(s.db, id); err != nil {
		s.logger.Error("删除操作日志失败", zap.Error(err))
		return fmt.Errorf("删除操作日志失败: %w", err)
	}

	s.logger.Info("删除操作日志成功", zap.Int64("id", id))
	return nil
}

// BatchDelete 批量删除操作日志
func (s *operLogService) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("日志ID列表不能为空")
	}

	// 批量删除
	rowsAffected, err := (&model.OperLog{}).BatchDelete(s.db, ids)
	if err != nil {
		s.logger.Error("批量删除操作日志失败", zap.Error(err))
		return fmt.Errorf("批量删除操作日志失败: %w", err)
	}

	s.logger.Info("批量删除操作日志成功", zap.Int64("count", rowsAffected))
	return nil
}

// GetById 根据ID查询操作日志
func (s *operLogService) GetById(ctx context.Context, id int64) (*model.OperLog, error) {
	log, err := (&model.OperLog{}).FindByID(s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("操作日志不存在")
		}
		s.logger.Error("查询操作日志失败", zap.Error(err))
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}
	return log, nil
}

// Page 分页查询操作日志列表
func (s *operLogService) Page(ctx context.Context, req *PageOperLogRequest) (*pagination.Page[model.OperLog], error) {
	// 1. 构建查询条件
	query := s.db.Model(&model.OperLog{})

	// 2. 添加条件过滤
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.OperName != "" {
		query = query.Where("oper_name LIKE ?", "%"+req.OperName+"%")
	}
	if req.BusinessType != "" {
		query = query.Where("business_type = ?", req.BusinessType)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartTime != "" {
		query = query.Where("oper_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("oper_time <= ?", req.EndTime)
	}

	// 3. 添加默认排序（如果 PageQuery 没有指定排序）
	if req.PageQuery.OrderByColumn == "" {
		query = query.Order("oper_time DESC, id DESC")
	}

	// 4. 使用 Paginator 执行分页查询
	page, err := pagination.New[model.OperLog](query, &req.PageQuery).Find()
	if err != nil {
		statusValue := ""
		if req.Status != nil {
			statusValue = *req.Status
		}
		s.logger.Error("查询操作日志列表失败",
			zap.Error(err),
			zap.String("title", req.Title),
			zap.String("operName", req.OperName),
			zap.String("businessType", req.BusinessType),
			zap.String("status", statusValue))
		return nil, fmt.Errorf("查询操作日志列表失败: %w", err)
	}

	return page, nil
}

// CleanOldLogs 清理指定天数之前的日志
func (s *operLogService) CleanOldLogs(ctx context.Context, days int) (int64, error) {
	if days <= 0 {
		return 0, fmt.Errorf("天数必须大于0")
	}

	rowsAffected, err := (&model.OperLog{}).CleanOldLogs(s.db, days)
	if err != nil {
		s.logger.Error("清理操作日志失败", zap.Error(err), zap.Int("days", days))
		return 0, fmt.Errorf("清理操作日志失败: %w", err)
	}

	s.logger.Info("清理操作日志成功", zap.Int64("count", rowsAffected), zap.Int("days", days))
	return rowsAffected, nil
}
