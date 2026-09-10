package sys

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys/internal/domain/model"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/modules"
	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
	"github.com/gcc798/microservice-kit/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ConfigService interface {
	Create(ctx context.Context, req *CreateConfigRequest) error
	Update(ctx context.Context, req *UpdateConfigRequest) error
	Delete(ctx context.Context, id int64) error
	BatchDelete(ctx context.Context, ids []int64) error
	GetById(ctx context.Context, id int64) (*model.Config, error)
	Page(ctx context.Context, pageNum, pageSize int, configCode, name string) (*pagination.Page[model.Config], error)
	GetByCode(ctx context.Context, configCode string) (*model.Config, error)
	GetDataByCode(ctx context.Context, configCode string) (json.RawMessage, error)
}

type ModuleRefresher interface {
	GetModule(name string) modules.Module
	RefreshModule(ctx context.Context, name string, req modules.ModuleRefreshRequest) error
}

type configService struct {
	db      *gorm.DB
	logger  logging.Logger
	store   *runtimeconfig.Store
	modules ModuleRefresher
}

func NewConfigService(db *gorm.DB, logger logging.Logger, store *runtimeconfig.Store, moduleRegistry ModuleRefresher) ConfigService {
	return &configService{db: db, logger: logger, store: store, modules: moduleRegistry}
}

func (s *configService) Create(ctx context.Context, req *CreateConfigRequest) error {
	if err := runtimeconfig.Validate(req.Code, req.Data); err != nil {
		return fmt.Errorf("配置数据无效: %w", err)
	}
	return s.store.WithCodeLock(ctx, req.Code, func() error {
		exists, err := (&model.Config{}).CheckNameExists(s.db.WithContext(ctx), req.Name)
		if err != nil {
			return fmt.Errorf("检查配置名称失败: %w", err)
		}
		if exists {
			return fmt.Errorf("配置名称已存在: %s", req.Name)
		}
		var count int64
		if err := s.db.WithContext(ctx).Model(&model.Config{}).Where("code = ?", req.Code).Count(&count).Error; err != nil {
			return fmt.Errorf("检查配置编码失败: %w", err)
		}
		if count != 0 {
			return fmt.Errorf("配置编码已存在: %s", req.Code)
		}
		if err := s.store.DeleteCache(ctx, req.Code); err != nil {
			return err
		}
		config := &model.Config{
			Name: req.Name, Code: req.Code, Data: req.Data, Remark: req.Remark,
			CreateBy: req.CreateBy, UpdateBy: req.UpdateBy,
		}
		if err := config.Create(s.db.WithContext(ctx)); err != nil {
			return fmt.Errorf("创建配置失败: %w", err)
		}
		if err := s.store.SetCache(ctx, req.Code, req.Data); err != nil {
			return err
		}
		s.logger.Info("创建配置成功", zap.Int64("id", config.ID), zap.String("code", config.Code))
		return s.refresh(ctx, req.Code, "created")
	})
}

func (s *configService) Update(ctx context.Context, req *UpdateConfigRequest) error {
	if err := runtimeconfig.Validate(req.Code, req.Data); err != nil {
		return fmt.Errorf("配置数据无效: %w", err)
	}
	existing, err := (&model.Config{}).FindByID(s.db.WithContext(ctx), req.ID)
	if err != nil {
		return configLookupError(err)
	}
	if req.Code != existing.Code {
		return errors.New("配置编码不可修改")
	}
	return s.store.WithCodeLock(ctx, existing.Code, func() error {
		current, err := (&model.Config{}).FindByID(s.db.WithContext(ctx), req.ID)
		if err != nil {
			return configLookupError(err)
		}
		if current.Code != existing.Code {
			return errors.New("配置编码不可修改")
		}
		if req.Name != current.Name {
			exists, err := (&model.Config{}).CheckNameExistsExcludingSelf(s.db.WithContext(ctx), req.ID, req.Name)
			if err != nil {
				return fmt.Errorf("检查配置名称失败: %w", err)
			}
			if exists {
				return fmt.Errorf("配置名称已被占用: %s", req.Name)
			}
		}
		if err := s.store.DeleteCache(ctx, current.Code); err != nil {
			return err
		}
		updates := map[string]any{
			"name": req.Name, "data": req.Data, "remark": req.Remark, "update_by": req.UpdateBy,
		}
		if err := current.Update(s.db.WithContext(ctx), req.ID, updates); err != nil {
			return fmt.Errorf("更新配置失败: %w", err)
		}
		if err := s.store.SetCache(ctx, current.Code, req.Data); err != nil {
			return err
		}
		s.logger.Info("更新配置成功", zap.Int64("id", req.ID), zap.String("code", current.Code))
		return s.refresh(ctx, current.Code, "updated")
	})
}

func (s *configService) Delete(ctx context.Context, id int64) error {
	existing, err := (&model.Config{}).FindByID(s.db.WithContext(ctx), id)
	if err != nil {
		return configLookupError(err)
	}
	return s.store.WithCodeLock(ctx, existing.Code, func() error {
		current, err := (&model.Config{}).FindByID(s.db.WithContext(ctx), id)
		if err != nil {
			return configLookupError(err)
		}
		if current.Code != existing.Code {
			return errors.New("配置编码在删除期间发生变化")
		}
		if err := s.store.DeleteCache(ctx, current.Code); err != nil {
			return err
		}
		if err := current.Delete(s.db.WithContext(ctx), id); err != nil {
			return fmt.Errorf("删除配置失败: %w", err)
		}
		s.logger.Info("删除配置成功", zap.Int64("id", id), zap.String("code", current.Code))
		return nil
	})
}

func (s *configService) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("配置ID列表不能为空")
	}
	for _, id := range ids {
		if err := s.Delete(ctx, id); err != nil {
			return fmt.Errorf("批量删除配置失败: %w", err)
		}
	}
	s.logger.Info("批量删除配置成功", zap.Int("count", len(ids)))
	return nil
}

func (s *configService) GetById(ctx context.Context, id int64) (*model.Config, error) {
	config, err := (&model.Config{}).FindByID(s.db.WithContext(ctx), id)
	if err != nil {
		return nil, configLookupError(err)
	}
	return config, nil
}

func (s *configService) Page(ctx context.Context, pageNum, pageSize int, configCode, name string) (*pagination.Page[model.Config], error) {
	query := s.db.WithContext(ctx).Model(&model.Config{})
	if configCode != "" {
		query = query.Where("code = ?", configCode)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	page, err := pagination.New[model.Config](query, &pagination.PageQuery{PageNum: pageNum, PageSize: pageSize}).Find()
	if err != nil {
		return nil, fmt.Errorf("分页查询配置列表失败: %w", err)
	}
	return page, nil
}

func (s *configService) GetByCode(ctx context.Context, configCode string) (*model.Config, error) {
	config, err := (&model.Config{}).FindByCode(s.db.WithContext(ctx), configCode)
	if err != nil {
		return nil, configLookupError(err)
	}
	return config, nil
}

func (s *configService) GetDataByCode(ctx context.Context, configCode string) (json.RawMessage, error) {
	data, err := (&model.Config{}).GetDataByCode(s.db.WithContext(ctx), configCode)
	if err != nil {
		return nil, configLookupError(err)
	}
	return data, nil
}

func (s *configService) refresh(ctx context.Context, code, reason string) error {
	name := moduleNameForCode(code)
	if name == "" || s.modules == nil || s.modules.GetModule(name) == nil {
		return nil
	}
	return s.modules.RefreshModule(ctx, name, modules.ModuleRefreshRequest{Codes: []string{code}, Reason: reason})
}

func moduleNameForCode(code string) string {
	switch code {
	case runtimeconfig.CodeWeChat:
		return modules.WeChatName
	case runtimeconfig.CodeSMS:
		return modules.SMSName
	case runtimeconfig.CodeEmail:
		return modules.EmailName
	case runtimeconfig.CodeCaptcha:
		return modules.CaptchaName
	case runtimeconfig.CodeScheduler:
		return modules.SchedulerName
	default:
		return ""
	}
}

func configLookupError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("配置不存在")
	}
	return fmt.Errorf("查询配置失败: %w", err)
}
