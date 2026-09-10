package sys

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys/internal/domain/model"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/utils"
	"github.com/gcc798/microservice-kit/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// LoginLogService 登录日志服务接口
type LoginLogService interface {
	// Create 创建登录日志
	Create(ctx context.Context, req *CreateLoginLogRequest) error

	// Update 更新登录日志
	Update(ctx context.Context, req *UpdateLoginLogRequest) error

	// Delete 删除登录日志
	Delete(ctx context.Context, id int64) error

	// BatchDelete 批量删除登录日志
	BatchDelete(ctx context.Context, ids []int64) error

	// GetById 根据ID查询登录日志
	GetById(ctx context.Context, id int64) (*model.LoginLog, error)

	// Page 分页查询登录日志列表
	Page(ctx context.Context, req *PageLoginLogRequest) (*pagination.Page[model.LoginLog], error)

	// CleanOldLogs 清理指定天数之前的日志
	CleanOldLogs(ctx context.Context, days int) (int64, error)
}

type loginLogService struct {
	db     *gorm.DB
	logger logging.Logger
}

// NewLoginLogService 创建登录日志服务实例
func NewLoginLogService(db *gorm.DB, logger logging.Logger) LoginLogService {
	return &loginLogService{
		db:     db,
		logger: logger,
	}
}

// Create 创建登录日志
func (s *loginLogService) Create(ctx context.Context, req *CreateLoginLogRequest) error {
	log := &model.LoginLog{
		UserName:      req.UserName,
		Ipaddr:        req.Ipaddr,
		LoginLocation: req.LoginLocation,
		Browser:       req.Browser,
		Os:            req.Os,
		Status:        req.Status,
		Msg:           req.Msg,
		LoginTime:     utils.Now(),
		ClientId:      req.ClientId,
	}

	if err := log.Create(s.db); err != nil {
		s.logger.Error("创建登录日志失败", zap.Error(err))
		return fmt.Errorf("创建登录日志失败: %w", err)
	}

	s.logger.Info("创建登录日志成功",
		zap.Int64("id", log.ID),
		zap.String("userName", log.UserName),
		zap.Int32("status", log.Status))

	return nil
}

// Update 更新登录日志
func (s *loginLogService) Update(ctx context.Context, req *UpdateLoginLogRequest) error {
	// 检查日志是否存在
	_, err := (&model.LoginLog{}).FindByID(s.db, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("登录日志不存在")
		}
		s.logger.Error("查询登录日志失败", zap.Error(err))
		return fmt.Errorf("查询登录日志失败: %w", err)
	}

	// 更新日志
	updates := map[string]interface{}{
		"user_name":      req.UserName,
		"ipaddr":         req.Ipaddr,
		"login_location": req.LoginLocation,
		"browser":        req.Browser,
		"os":             req.Os,
		"status":         req.Status,
		"msg":            req.Msg,
		"client_id":      req.ClientId,
	}

	log := &model.LoginLog{}
	if err := log.Update(s.db, req.ID, updates); err != nil {
		s.logger.Error("更新登录日志失败", zap.Error(err))
		return fmt.Errorf("更新登录日志失败: %w", err)
	}

	s.logger.Info("更新登录日志成功", zap.Int64("id", req.ID))
	return nil
}

// Delete 删除登录日志
func (s *loginLogService) Delete(ctx context.Context, id int64) error {
	// 检查日志是否存在
	_, err := (&model.LoginLog{}).FindByID(s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("登录日志不存在")
		}
		s.logger.Error("查询登录日志失败", zap.Error(err))
		return fmt.Errorf("查询登录日志失败: %w", err)
	}

	// 删除日志
	if err := (&model.LoginLog{}).Delete(s.db, id); err != nil {
		s.logger.Error("删除登录日志失败", zap.Error(err))
		return fmt.Errorf("删除登录日志失败: %w", err)
	}

	s.logger.Info("删除登录日志成功", zap.Int64("id", id))
	return nil
}

// BatchDelete 批量删除登录日志
func (s *loginLogService) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("日志ID列表不能为空")
	}

	// 批量删除
	rowsAffected, err := (&model.LoginLog{}).BatchDelete(s.db, ids)
	if err != nil {
		s.logger.Error("批量删除登录日志失败", zap.Error(err))
		return fmt.Errorf("批量删除登录日志失败: %w", err)
	}

	s.logger.Info("批量删除登录日志成功", zap.Int64("count", rowsAffected))
	return nil
}

// GetById 根据ID查询登录日志
func (s *loginLogService) GetById(ctx context.Context, id int64) (*model.LoginLog, error) {
	log, err := (&model.LoginLog{}).FindByID(s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("登录日志不存在")
		}
		s.logger.Error("查询登录日志失败", zap.Error(err))
		return nil, fmt.Errorf("查询登录日志失败: %w", err)
	}
	return log, nil
}

// Page 分页查询登录日志列表
func (s *loginLogService) Page(ctx context.Context, req *PageLoginLogRequest) (*pagination.Page[model.LoginLog], error) {
	// 1. 构建查询条件
	query := s.db.Model(&model.LoginLog{})

	// 2. 添加条件过滤
	if req.UserName != "" {
		query = query.Where("user_name LIKE ?", "%"+req.UserName+"%")
	}
	if req.Ipaddr != "" {
		query = query.Where("ipaddr LIKE ?", "%"+req.Ipaddr+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartTime != "" {
		query = query.Where("login_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("login_time <= ?", req.EndTime)
	}

	// 3. 添加默认排序（如果 PageQuery 没有指定排序）
	if req.PageQuery.OrderByColumn == "" {
		query = query.Order("login_time DESC, id DESC")
	}

	// 4. 使用 Paginator 执行分页查询
	page, err := pagination.New[model.LoginLog](query, &req.PageQuery).Find()
	if err != nil {
		statusValue := int32(-1)
		if req.Status != nil {
			statusValue = *req.Status
		}
		s.logger.Error("查询登录日志列表失败",
			zap.Error(err),
			zap.String("userName", req.UserName),
			zap.String("ipaddr", req.Ipaddr),
			zap.Int32("status", statusValue))
		return nil, fmt.Errorf("查询登录日志列表失败: %w", err)
	}

	return page, nil
}

// CleanOldLogs 清理指定天数之前的日志
func (s *loginLogService) CleanOldLogs(ctx context.Context, days int) (int64, error) {
	if days <= 0 {
		return 0, fmt.Errorf("天数必须大于0")
	}

	rowsAffected, err := (&model.LoginLog{}).CleanOldLogs(s.db, days)
	if err != nil {
		s.logger.Error("清理登录日志失败", zap.Error(err), zap.Int("days", days))
		return 0, fmt.Errorf("清理登录日志失败: %w", err)
	}

	s.logger.Info("清理登录日志成功", zap.Int64("count", rowsAffected), zap.Int("days", days))
	return rowsAffected, nil
}
