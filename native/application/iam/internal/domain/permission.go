package iam

import (
	"context"
	"fmt"

	logging "github.com/gcc798/lightning/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PermissionService 基于业务表检查用户的 API 权限。
// 有效权限由用户直接权限和用户启用角色的权限共同组成。
type PermissionService interface {
	CheckPermission(ctx context.Context, userId int64, resource, action string) (bool, error)
}

type permissionService struct {
	db     *gorm.DB
	logger logging.Logger
}

// NewPermissionService 创建数据库权限服务。
func NewPermissionService(db *gorm.DB, logger logging.Logger) PermissionService {
	return &permissionService{db: db, logger: logger}
}

// CheckPermission 检查超级管理员角色、角色权限和用户直接权限。
func (s *permissionService) CheckPermission(ctx context.Context, userId int64, resource, action string) (bool, error) {
	var allowed bool
	err := s.db.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM m_user_role ur
			JOIN s_role r ON r.id = ur.role_id
			WHERE ur.user_id = ?
			  AND r.status = 0
			  AND r.role_key = 'super_admin'

			UNION ALL

			SELECT 1
			FROM m_user_role ur
			JOIN s_role r ON r.id = ur.role_id
			JOIN m_role_api_permission rp ON rp.role_id = r.id
			JOIN s_api_permission p ON p.id = rp.permission_id
			WHERE ur.user_id = ?
			  AND r.status = 0
			  AND p.status = 0
			  AND (p.action = ? OR p.action = '*')
			  AND (
				p.code = ?
				OR p.code = '*'
				OR (RIGHT(p.code, 1) = '*' AND ? LIKE LEFT(p.code, LENGTH(p.code) - 1) || '%')
				OR (LEFT(p.code, 1) = '*' AND ? LIKE '%' || SUBSTRING(p.code FROM 2))
			  )

			UNION ALL

			SELECT 1
			FROM m_user_api_permission up
			JOIN s_api_permission p ON p.id = up.permission_id
			WHERE up.user_id = ?
			  AND p.status = 0
			  AND (p.action = ? OR p.action = '*')
			  AND (
				p.code = ?
				OR p.code = '*'
				OR (RIGHT(p.code, 1) = '*' AND ? LIKE LEFT(p.code, LENGTH(p.code) - 1) || '%')
				OR (LEFT(p.code, 1) = '*' AND ? LIKE '%' || SUBSTRING(p.code FROM 2))
			  )
		) AS allowed
	`,
		userId,
		userId, action, resource, resource, resource,
		userId, action, resource, resource, resource,
	).Scan(&allowed).Error
	if err != nil {
		s.logger.Error("权限检查失败",
			zap.Int64("userId", userId),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Error(err))
		return false, fmt.Errorf("权限检查失败: %w", err)
	}

	s.logger.Debug("权限检查",
		zap.Int64("userId", userId),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.Bool("allowed", allowed))
	return allowed, nil
}
