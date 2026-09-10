package data

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent/casbinrule"
)

func (r *Resources) activeUserRoles(ctx context.Context) ([]*entpkg.UserRole, error) {
	items, err := r.Ent.UserRole.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.UserRole, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *Resources) activeRoleMenus(ctx context.Context) ([]*entpkg.RoleMenu, error) {
	items, err := r.Ent.RoleMenu.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.RoleMenu, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *Resources) AddUserRole(ctx context.Context, userID, roleID int64) error {
	user, err := r.Ent.User.Get(ctx, userID)
	if err != nil {
		return err
	}
	if user.DeletedAt != nil {
		return errors.New("user not found")
	}
	role, err := r.Ent.Role.Get(ctx, roleID)
	if err != nil {
		return err
	}
	if role.DeletedAt != nil {
		return errors.New("role not found")
	}
	items, err := r.activeUserRoles(ctx)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.UserID == userID && item.RoleID == roleID {
			return r.EnsureUserRoleRule(ctx, userID, role.RoleKey)
		}
	}
	now := time.Now()
	operator := currentOperatorID(ctx)
	_, err = r.Ent.UserRole.Create().SetID(nextID()).SetUserID(userID).SetRoleID(roleID).SetCreateBy(operator).SetUpdateBy(operator).SetCreatedTime(now).SetUpdatedTime(now).Save(ctx)
	if err != nil {
		return err
	}
	return r.EnsureUserRoleRule(ctx, userID, role.RoleKey)
}

func (r *Resources) RemoveUserRole(ctx context.Context, userID, roleID int64) error {
	if _, err := r.Ent.User.Get(ctx, userID); err != nil {
		return err
	}
	role, err := r.Ent.Role.Get(ctx, roleID)
	if err != nil {
		return err
	}
	items, err := r.activeUserRoles(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	operator := currentOperatorID(ctx)
	for _, item := range items {
		if item.UserID == userID && item.RoleID == roleID {
			if _, err = r.Ent.UserRole.UpdateOneID(item.ID).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil {
				return err
			}
			return r.RemoveUserRoleRule(ctx, userID, role.RoleKey)
		}
	}
	return nil
}

func (r *Resources) GetUserRoles(ctx context.Context, userID int64) ([]*v1.RoleItem, error) {
	user, err := r.Ent.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.DeletedAt != nil {
		return nil, errors.New("user not found")
	}
	userRoles, err := r.activeUserRoles(ctx)
	if err != nil {
		return nil, err
	}
	roleMap := make(map[int64]struct{})
	for _, item := range userRoles {
		if item.UserID == userID {
			roleMap[item.RoleID] = struct{}{}
		}
	}
	roles, err := r.activeRoles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.RoleItem, 0)
	for _, item := range roles {
		if _, ok := roleMap[item.ID]; ok {
			out = append(out, roleEntityToItem(item))
		}
	}
	return out, nil
}

func (r *Resources) AddRolePermission(ctx context.Context, roleKey, resource, action string) error {
	roleKey = strings.TrimSpace(roleKey)
	resource = strings.TrimSpace(resource)
	action = strings.TrimSpace(action)
	if roleKey == "" || resource == "" {
		return errors.New("roleKey和resource不能为空")
	}
	if action == "" {
		action = permissionAction(resource)
	}
	exists, err := r.Ent.CasbinRule.Query().
		Where(
			casbinrule.Ptype("p"),
			casbinrule.V0(fmt.Sprintf("role::%s", roleKey)),
			casbinrule.V1(resource),
			casbinrule.V2(action),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = r.Ent.CasbinRule.Create().
		SetID(nextID()).
		SetPtype("p").
		SetV0(fmt.Sprintf("role::%s", roleKey)).
		SetV1(resource).
		SetV2(action).
		Save(ctx)
	return err
}

func (r *Resources) DeleteRolePermission(ctx context.Context, roleKey, resource, action string) error {
	roleKey = strings.TrimSpace(roleKey)
	resource = strings.TrimSpace(resource)
	action = strings.TrimSpace(action)
	if roleKey == "" || resource == "" {
		return nil
	}
	if action == "" {
		action = permissionAction(resource)
	}
	_, err := r.Ent.CasbinRule.Delete().
		Where(
			casbinrule.Ptype("p"),
			casbinrule.V0(fmt.Sprintf("role::%s", roleKey)),
			casbinrule.V1(resource),
			casbinrule.V2(action),
		).
		Exec(ctx)
	return err
}

func (r *Resources) GetRolePermissions(ctx context.Context, roleKey string) ([]string, error) {
	roleKey = strings.TrimSpace(roleKey)
	if roleKey == "" {
		return nil, nil
	}
	items, err := r.Ent.CasbinRule.Query().
		Where(casbinrule.Ptype("p"), casbinrule.V0(fmt.Sprintf("role::%s", roleKey))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	perms := make([]string, 0, len(items))
	for _, item := range items {
		if item.V1 == "" {
			continue
		}
		perms = append(perms, joinResourceAction(item.V1, item.V2))
	}
	sort.Strings(perms)
	return perms, nil
}

func (r *Resources) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	if userID == 1 {
		return true, nil
	}
	if userID <= 0 || strings.TrimSpace(resource) == "" {
		return false, nil
	}
	action = strings.TrimSpace(action)
	if action == "" {
		action = permissionAction(resource)
	}
	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}
	if len(roles) == 0 {
		return false, nil
	}
	subjects := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.GetRoleKey() != "" {
			subjects = append(subjects, fmt.Sprintf("role::%s", role.GetRoleKey()))
		}
	}
	if len(subjects) == 0 {
		return false, nil
	}
	rules, err := r.Ent.CasbinRule.Query().
		Where(casbinrule.Ptype("p"), casbinrule.V0In(subjects...)).
		All(ctx)
	if err != nil {
		return false, err
	}
	for _, rule := range rules {
		if resourceMatch(rule.V1, resource) && actionMatch(rule.V2, action) {
			return true, nil
		}
	}
	return false, nil
}

func permissionAction(resource string) string {
	resource = strings.TrimSpace(resource)
	if strings.HasSuffix(resource, ".read") {
		return "read"
	}
	return "write"
}

func joinResourceAction(resource, action string) string {
	resource = strings.TrimSpace(resource)
	action = strings.TrimSpace(action)
	if resource == "" {
		return ""
	}
	if action == "" {
		return resource
	}
	return fmt.Sprintf("%s:%s", resource, action)
}

func resourceMatch(policy, resource string) bool {
	policy = strings.TrimSpace(policy)
	resource = strings.TrimSpace(resource)
	switch {
	case policy == "" || resource == "":
		return false
	case policy == "*" || policy == resource:
		return true
	case strings.HasSuffix(policy, ".*"):
		return strings.HasPrefix(resource, strings.TrimSuffix(policy, ".*")+".")
	case strings.HasPrefix(policy, "*."):
		return strings.HasSuffix(resource, strings.TrimPrefix(policy, "*"))
	default:
		return false
	}
}

func actionMatch(policy, action string) bool {
	policy = strings.TrimSpace(policy)
	action = strings.TrimSpace(action)
	return policy == "*" || policy == action
}
