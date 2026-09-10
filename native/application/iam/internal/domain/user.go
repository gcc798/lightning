package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/microservice-kit/application/iam/internal/domain/model"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/utils"
	apperrors "github.com/gcc798/microservice-kit/internal/utils/errors"
	"github.com/gcc798/microservice-kit/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	// Create 创建用户
	Create(ctx context.Context, req *CreateUserRequest) error

	// Update 更新用户
	Update(ctx context.Context, req *UpdateUserRequest) error

	// Delete 删除单个用户
	Delete(ctx context.Context, userId int64) error

	// BatchDelete 批量删除用户
	BatchDelete(ctx context.Context, userIds []int64) error

	// GetById 根据ID查询用户
	GetById(ctx context.Context, userId int64) (*model.User, error)

	// Page 分页查询用户列表
	Page(ctx context.Context, pageNum, pageSize int, username, phonenumber string, status int32) (*pagination.Page[model.User], error)

	// BatchImport 批量导入用户
	BatchImport(ctx context.Context, req *BatchImportUsersRequest) (successCount int, failCount int, errors []string, err error)

	// ResetPassword 重置用户密码
	ResetPassword(ctx context.Context, userId int64, newPassword string) error

	// ChangePassword 用户修改密码
	ChangePassword(ctx context.Context, userId int64, oldPassword, newPassword string) error
}

type userService struct {
	db     *gorm.DB
	logger logging.Logger
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB, logger logging.Logger) UserService {
	return &userService{
		db:     db,
		logger: logger,
	}
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req *CreateUserRequest) error {
	// 一次查询检查所有冲突（用户名、手机号、邮箱）
	conflicts, err := (&model.User{}).FindConflicts(s.db, req.UserName, req.Phonenumber, req.Email)
	if err != nil {
		s.logger.Error("检查冲突失败", zap.Error(err))
		return fmt.Errorf("检查冲突失败: %w", err)
	}

	// 在内存中进行多次校验，避免重复查询数据库
	for _, user := range conflicts {
		if user.UserName == req.UserName {
			return apperrors.NewBusiness(apperrors.CodeUserNameExists, "用户名已存在")
		}
		if req.Phonenumber != "" && user.Phonenumber == req.Phonenumber {
			return apperrors.NewBusiness(apperrors.CodePhoneExists, "手机号已存在")
		}
		if req.Email != "" && user.Email == req.Email {
			return apperrors.NewBusiness(apperrors.CodeEmailExists, "邮箱已存在")
		}
	}

	// 加密密码
	var hashedPassword string
	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			s.logger.Error("密码加密失败", zap.Error(err))
			return fmt.Errorf("密码加密失败: %w", err)
		}
		hashedPassword = hashed
	}

	// 创建用户实体
	user := &model.User{
		UserName:    req.UserName,
		NickName:    req.NickName,
		Password:    hashedPassword,
		UserType:    req.UserType,
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		Sex:         req.Sex,
		Avatar:      req.Avatar,
		Status:      req.Status,
		Remark:      req.Remark,
		CreateBy:    req.CreateBy,
		UpdateBy:    req.UpdateBy,
	}

	// 调用模型层的创建方法
	if err := user.Create(s.db, user); err != nil {
		s.logger.Error("创建用户失败", zap.Error(err))
		return fmt.Errorf("创建用户失败: %w", err)
	}

	s.logger.Info("创建用户成功", zap.Int64("userId", user.ID), zap.String("userName", user.UserName))
	return nil
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, req *UpdateUserRequest) error {
	// 检查用户是否存在
	existingUser, err := (&model.User{}).FindByID(s.db, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewBusiness(apperrors.CodeUserNotFound, "用户不存在")
		}
		s.logger.Error("查询用户失败", zap.Error(err))
		return apperrors.NewInfrastructure(apperrors.CodeDatabaseError, "数据库查询失败", err)
	}

	// 一次查询检查所有冲突（排除自己）
	conflicts, err := (&model.User{}).FindConflictsExcludingSelf(
		s.db, req.UserId, req.UserName, req.Phonenumber, req.Email,
	)
	if err != nil {
		s.logger.Error("检查冲突失败", zap.Error(err))
		return fmt.Errorf("检查冲突失败: %w", err)
	}

	// 在内存中进行多次校验
	for _, user := range conflicts {
		if req.UserName != "" && user.UserName == req.UserName {
			return apperrors.NewBusiness(apperrors.CodeUserNameExists, "用户名已被占用")
		}
		if req.Phonenumber != "" && user.Phonenumber == req.Phonenumber {
			return apperrors.NewBusiness(apperrors.CodePhoneExists, "手机号已被占用")
		}
		if req.Email != "" && user.Email == req.Email {
			return apperrors.NewBusiness(apperrors.CodeEmailExists, "邮箱已被占用")
		}
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"user_name":   req.UserName,
		"nick_name":   req.NickName,
		"user_type":   req.UserType,
		"email":       req.Email,
		"phonenumber": req.Phonenumber,
		"sex":         req.Sex,
		"avatar":      req.Avatar,
		"status":      req.Status,
		"remark":      req.Remark,
		"update_by":   req.UpdateBy,
	}

	// 调用模型层的更新方法
	if err := existingUser.Update(s.db, req.UserId, updates); err != nil {
		s.logger.Error("更新用户失败", zap.Error(err))
		return fmt.Errorf("更新用户失败: %w", err)
	}

	s.logger.Info("更新用户成功", zap.Int64("userId", req.UserId))
	return nil
}

// Delete 删除单个用户
func (s *userService) Delete(ctx context.Context, userId int64) error {
	// 检查用户是否存在
	user, err := (&model.User{}).FindByID(s.db, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewBusiness(apperrors.CodeUserNotFound, "用户不存在")
		}
		s.logger.Error("查询用户失败", zap.Error(err))
		return apperrors.NewInfrastructure(apperrors.CodeDatabaseError, "数据库查询失败", err)
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userId).Delete(&model.MUserRole{}).Error; err != nil {
			return fmt.Errorf("删除用户角色关联失败: %w", err)
		}
		if err := tx.Where("user_id = ?", userId).Delete(&model.MUserApiPermission{}).Error; err != nil {
			return fmt.Errorf("删除用户 API 权限关联失败: %w", err)
		}
		return user.Delete(tx, userId)
	}); err != nil {
		s.logger.Error("删除用户失败", zap.Error(err))
		return fmt.Errorf("删除用户失败: %w", err)
	}

	s.logger.Info("删除用户成功", zap.Int64("userId", userId))
	return nil
}

// BatchDelete 批量删除用户
func (s *userService) BatchDelete(ctx context.Context, userIds []int64) error {
	if len(userIds) == 0 {
		return fmt.Errorf("用户ID列表不能为空")
	}

	var rowsAffected int64
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id IN ?", userIds).Delete(&model.MUserRole{}).Error; err != nil {
			return fmt.Errorf("删除用户角色关联失败: %w", err)
		}
		if err := tx.Where("user_id IN ?", userIds).Delete(&model.MUserApiPermission{}).Error; err != nil {
			return fmt.Errorf("删除用户 API 权限关联失败: %w", err)
		}
		var err error
		rowsAffected, err = (&model.User{}).BatchDelete(tx, userIds)
		return err
	})
	if err != nil {
		s.logger.Error("批量删除用户失败", zap.Error(err))
		return fmt.Errorf("批量删除用户失败: %w", err)
	}

	s.logger.Info("批量删除用户成功", zap.Int64("count", rowsAffected))
	return nil
}

// GetById 根据ID查询用户
func (s *userService) GetById(ctx context.Context, userId int64) (*model.User, error) {
	user, err := (&model.User{}).FindByID(s.db, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBusiness(apperrors.CodeUserNotFound, "用户不存在")
		}
		s.logger.Error("查询用户失败", zap.Error(err))
		return nil, apperrors.NewInfrastructure(apperrors.CodeDatabaseError, "数据库查询失败", err)
	}

	// 清空密码字段
	user.ClearPassword()
	return user, nil
}

// Page 分页查询用户列表
func (s *userService) Page(ctx context.Context, pageNum, pageSize int, username, phonenumber string, status int32) (*pagination.Page[model.User], error) {
	query := s.db.Model(&model.User{})

	// 条件查询
	if username != "" {
		query = query.Where("user_name LIKE ?", "%"+username+"%")
	}
	if phonenumber != "" {
		query = query.Where("phonenumber LIKE ?", "%"+phonenumber+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 构建 PageQuery
	pageQuery := &pagination.PageQuery{
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	// 使用 Paginator 进行分页
	page, err := pagination.New[model.User](query, pageQuery).Find()
	if err != nil {
		s.logger.Error("分页查询用户列表失败", zap.Error(err))
		return nil, fmt.Errorf("分页查询用户列表失败: %w", err)
	}

	// 清空密码字段
	for i := range page.Records {
		page.Records[i].ClearPassword()
	}

	return page, nil
}

// BatchImport 批量导入用户
func (s *userService) BatchImport(ctx context.Context, req *BatchImportUsersRequest) (int, int, []string, error) {
	if len(req.Users) == 0 {
		return 0, 0, nil, fmt.Errorf("导入用户列表不能为空")
	}

	successCount := 0
	failCount := 0
	var errors []string

	// 逐个导入用户
	for i, userReq := range req.Users {
		// 加密密码
		var hashedPassword string
		if userReq.Password != "" {
			hashed, err := utils.HashPassword(userReq.Password)
			if err != nil {
				failCount++
				errors = append(errors, fmt.Sprintf("第%d行: 密码加密失败", i+1))
				continue
			}
			hashedPassword = hashed
		}

		// 创建用户
		user := &model.User{
			UserName:    userReq.UserName,
			NickName:    userReq.NickName,
			Password:    hashedPassword,
			UserType:    userReq.UserType,
			Email:       userReq.Email,
			Phonenumber: userReq.Phonenumber,
			Sex:         userReq.Sex,
			Avatar:      userReq.Avatar,
			Status:      userReq.Status,
			Remark:      userReq.Remark,
			CreateBy:    userReq.CreateBy,
			UpdateBy:    userReq.UpdateBy,
		}

		if err := s.db.Create(user).Error; err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d行: %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	s.logger.Info("批量导入用户完成",
		zap.Int("successCount", successCount),
		zap.Int("failCount", failCount))

	return successCount, failCount, errors, nil
}

// ResetPassword 重置用户密码
func (s *userService) ResetPassword(ctx context.Context, userId int64, newPassword string) error {
	// 检查用户是否存在
	user, err := (&model.User{}).FindByID(s.db, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewBusiness(apperrors.CodeUserNotFound, "用户不存在")
		}
		s.logger.Error("查询用户失败", zap.Error(err))
		return apperrors.NewInfrastructure(apperrors.CodeDatabaseError, "数据库查询失败", err)
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		s.logger.Error("密码加密失败", zap.Error(err))
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 调用模型层的更新密码方法
	if err := user.UpdatePassword(s.db, userId, hashedPassword); err != nil {
		s.logger.Error("重置密码失败", zap.Error(err))
		return fmt.Errorf("重置密码失败: %w", err)
	}

	s.logger.Info("重置密码成功", zap.Int64("userId", userId))
	return nil
}

// ChangePassword 用户修改密码
func (s *userService) ChangePassword(ctx context.Context, userId int64, oldPassword, newPassword string) error {
	// 1. 检查用户是否存在
	user, err := (&model.User{}).FindByID(s.db, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewBusiness(apperrors.CodeUserNotFound, "用户不存在")
		}
		s.logger.Error("查询用户失败", zap.Error(err))
		return apperrors.NewInfrastructure(apperrors.CodeDatabaseError, "数据库查询失败", err)
	}

	// 2. 验证旧密码
	if err := utils.VerifyPassword(user.Password, oldPassword); err != nil {
		s.logger.Warn("旧密码验证失败", zap.Int64("userId", userId))
		return apperrors.NewBusiness(apperrors.CodeInvalidPassword, "旧密码不正确")
	}

	// 3. 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		s.logger.Error("密码加密失败", zap.Error(err))
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 4. 更新密码
	if err := user.UpdatePassword(s.db, userId, hashedPassword); err != nil {
		s.logger.Error("修改密码失败", zap.Error(err))
		return fmt.Errorf("修改密码失败: %w", err)
	}

	s.logger.Info("修改密码成功", zap.Int64("userId", userId))
	return nil
}
