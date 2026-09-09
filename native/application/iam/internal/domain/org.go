package iam

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gcc798/lightning/application/iam/internal/domain/model"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OrgService 组织服务接口
type OrgService interface {
	// Create 创建组织，返回组织ID
	Create(ctx context.Context, req *CreateOrgRequest) (int64, error)

	// Update 更新组织
	Update(ctx context.Context, req *UpdateOrgRequest) error

	// Delete 删除单个组织
	Delete(ctx context.Context, orgId int64) error

	// BatchDelete 批量删除组织
	BatchDelete(ctx context.Context, orgIds []int64) error

	// GetById 根据ID查询组织
	GetById(ctx context.Context, orgId int64) (*model.Org, error)

	// Page 分页查询组织列表
	Page(ctx context.Context, pageNum, pageSize int, orgName, orgCode string, status int32, parentId *int64) (*pagination.Page[model.Org], error)

	// GetTree 获取组织树
	GetTree(ctx context.Context) ([]*OrgTree, error)
}

type orgService struct {
	db     *gorm.DB
	logger logging.Logger
}

// NewOrgService 创建组织服务实例
func NewOrgService(db *gorm.DB, logger logging.Logger) OrgService {
	return &orgService{
		db:     db,
		logger: logger,
	}
}

// Create 创建组织，返回组织ID
func (s *orgService) Create(ctx context.Context, req *CreateOrgRequest) (int64, error) {
	// 检查组织编码唯一性
	exists, err := (&model.Org{}).CheckOrgCodeExists(s.db, req.OrgCode)
	if err != nil {
		s.logger.Error("检查组织编码失败", zap.Error(err))
		return 0, fmt.Errorf("检查组织编码失败: %w", err)
	}
	if exists {
		return 0, errors.New("组织编码已存在")
	}

	// 构建祖级列表
	ancestors, err := (&model.Org{}).BuildAncestors(s.db, req.ParentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("父组织不存在")
		}
		s.logger.Error("构建祖级列表失败", zap.Error(err))
		return 0, fmt.Errorf("构建祖级列表失败: %w", err)
	}

	// 创建组织对象
	org := &model.Org{
		ParentId:  req.ParentId,
		Ancestors: ancestors,
		OrgName:   req.OrgName,
		OrgCode:   req.OrgCode,
		OrgType:   req.OrgType,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
		Sort:      req.Sort,
		Remark:    req.Remark,
		CreateBy:  req.CreateBy,
		UpdateBy:  req.UpdateBy,
	}

	// 设置默认值
	if org.Status == 0 {
		org.Status = 0 // StatusNormal
	}
	if org.OrgType == "" {
		org.OrgType = "company"
	}

	// 调用模型层的创建方法
	if err := org.Create(s.db, org); err != nil {
		s.logger.Error("创建组织失败", zap.Error(err))
		return 0, fmt.Errorf("创建组织失败: %w", err)
	}

	s.logger.Info("创建组织成功", zap.Int64("orgId", org.ID), zap.String("orgName", org.OrgName))
	return org.ID, nil
}

// Update 更新组织
func (s *orgService) Update(ctx context.Context, req *UpdateOrgRequest) error {
	if req.OrgId == 0 {
		return errors.New("组织ID不能为空")
	}

	// 检查组织是否存在
	existingOrg, err := (&model.Org{}).FindByID(s.db, req.OrgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("组织不存在")
		}
		s.logger.Error("查询组织失败", zap.Error(err))
		return fmt.Errorf("查询组织失败: %w", err)
	}

	// 检查组织编码是否被其他组织占用
	if req.OrgCode != "" && req.OrgCode != existingOrg.OrgCode {
		exists, err := (&model.Org{}).CheckOrgCodeExistsExcludingSelf(s.db, req.OrgId, req.OrgCode)
		if err != nil {
			s.logger.Error("检查组织编码失败", zap.Error(err))
			return fmt.Errorf("检查组织编码失败: %w", err)
		}
		if exists {
			return errors.New("组织编码已被占用")
		}
	}

	// 检查是否修改了父组织
	if req.ParentId != existingOrg.ParentId {
		// 不能将组织设置为自己的子组织
		if req.ParentId == req.OrgId {
			return errors.New("不能将组织设置为自己的子组织")
		}

		// 重新构建祖级列表
		ancestors, err := (&model.Org{}).BuildAncestors(s.db, req.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父组织不存在")
			}
			s.logger.Error("构建祖级列表失败", zap.Error(err))
			return fmt.Errorf("构建祖级列表失败: %w", err)
		}

		// 检查父组织的祖级列表中是否包含当前组织（避免循环引用）
		if strings.Contains(ancestors, strconv.FormatInt(req.OrgId, 10)) {
			return errors.New("不能将组织移动到其子组织下")
		}

		existingOrg.Ancestors = ancestors
	}

	// 更新字段
	if req.ParentId != 0 {
		existingOrg.ParentId = req.ParentId
	}
	if req.OrgName != "" {
		existingOrg.OrgName = req.OrgName
	}
	if req.OrgCode != "" {
		existingOrg.OrgCode = req.OrgCode
	}
	if req.OrgType != "" {
		existingOrg.OrgType = req.OrgType
	}
	existingOrg.Leader = req.Leader
	existingOrg.Phone = req.Phone
	existingOrg.Email = req.Email
	if req.Status != 0 {
		existingOrg.Status = req.Status
	}
	existingOrg.Sort = req.Sort
	existingOrg.Remark = req.Remark
	existingOrg.UpdateBy = req.UpdateBy

	// 调用模型层的更新方法
	if err := existingOrg.Update(s.db, existingOrg); err != nil {
		s.logger.Error("更新组织失败", zap.Error(err))
		return fmt.Errorf("更新组织失败: %w", err)
	}

	s.logger.Info("更新组织成功", zap.Int64("orgId", req.OrgId))
	return nil
}

// Delete 删除单个组织
func (s *orgService) Delete(ctx context.Context, orgId int64) error {
	if orgId == 0 {
		return errors.New("组织ID不能为空")
	}

	// 检查组织是否存在
	org, err := (&model.Org{}).FindByID(s.db, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("组织不存在")
		}
		s.logger.Error("查询组织失败", zap.Error(err))
		return fmt.Errorf("查询组织失败: %w", err)
	}

	// 检查是否有子组织
	hasChildren, err := org.HasChildren(s.db)
	if err != nil {
		s.logger.Error("检查子组织失败", zap.Error(err))
		return fmt.Errorf("检查子组织失败: %w", err)
	}
	if hasChildren {
		return errors.New("存在子组织，无法删除")
	}

	// 检查是否有关联用户
	hasUsers, err := org.HasUsers(s.db)
	if err != nil {
		s.logger.Error("检查关联用户失败", zap.Error(err))
		return fmt.Errorf("检查关联用户失败: %w", err)
	}
	if hasUsers {
		return errors.New("组织下存在用户，无法删除")
	}

	// 调用模型层的删除方法
	if err := org.Delete(s.db, orgId); err != nil {
		s.logger.Error("删除组织失败", zap.Error(err))
		return fmt.Errorf("删除组织失败: %w", err)
	}

	s.logger.Info("删除组织成功", zap.Int64("orgId", orgId))
	return nil
}

// BatchDelete 批量删除组织
func (s *orgService) BatchDelete(ctx context.Context, orgIds []int64) error {
	if len(orgIds) == 0 {
		return errors.New("组织ID列表不能为空")
	}

	// 逐个删除（需要检查业务规则）
	var failedIds []int64
	var errors []string

	for _, orgId := range orgIds {
		if err := s.Delete(ctx, orgId); err != nil {
			failedIds = append(failedIds, orgId)
			errors = append(errors, fmt.Sprintf("组织ID %d: %s", orgId, err.Error()))
		}
	}

	if len(failedIds) > 0 {
		return fmt.Errorf("部分组织删除失败: %s", strings.Join(errors, "; "))
	}

	s.logger.Info("批量删除组织成功", zap.Int("count", len(orgIds)))
	return nil
}

// GetById 根据ID查询组织
func (s *orgService) GetById(ctx context.Context, orgId int64) (*model.Org, error) {
	if orgId == 0 {
		return nil, errors.New("组织ID不能为空")
	}

	org, err := (&model.Org{}).FindByID(s.db, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("组织不存在")
		}
		s.logger.Error("查询组织失败", zap.Error(err))
		return nil, fmt.Errorf("查询组织失败: %w", err)
	}

	return org, nil
}

// Page 分页查询组织列表
func (s *orgService) Page(ctx context.Context, pageNum, pageSize int, orgName, orgCode string, status int32, parentId *int64) (*pagination.Page[model.Org], error) {
	query := s.db.Model(&model.Org{})

	// 条件查询
	if orgName != "" {
		query = query.Where("org_name LIKE ?", "%"+orgName+"%")
	}
	if orgCode != "" {
		query = query.Where("org_code = ?", orgCode)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if parentId != nil {
		query = query.Where("parent_id = ?", *parentId)
	}

	// 构建 PageQuery
	pageQuery := &pagination.PageQuery{
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	// 使用 Paginator 进行分页
	page, err := pagination.New[model.Org](query, pageQuery).Find()
	if err != nil {
		s.logger.Error("分页查询组织列表失败", zap.Error(err))
		return nil, fmt.Errorf("分页查询组织列表失败: %w", err)
	}

	return page, nil
}

// OrgTree 组织树节点
type OrgTree struct {
	model.Org
	Children []*OrgTree `json:"children,omitempty"`
}

// GetTree 获取组织树（所有组织）
func (s *orgService) GetTree(ctx context.Context) ([]*OrgTree, error) {
	orgs, err := (&model.Org{}).FindAll(s.db)
	if err != nil {
		s.logger.Error("查询组织树失败", zap.Error(err))
		return nil, fmt.Errorf("查询组织树失败: %w", err)
	}
	return s.buildOrgTree(orgs, 0), nil
}

// buildOrgTree 构建组织树
func (s *orgService) buildOrgTree(orgs []model.Org, parentId int64) []*OrgTree {
	var tree []*OrgTree

	for _, org := range orgs {
		if org.ParentId == parentId {
			node := &OrgTree{
				Org:      org,
				Children: s.buildOrgTree(orgs, org.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
