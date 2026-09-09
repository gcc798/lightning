package iam

// ApiPermissionSaveRequest 保存 API 权限节点。
type ApiPermissionSaveRequest struct {
	ParentId int64  `json:"parentId"`
	Module   string `json:"module" binding:"required,max=64"`
	Code     string `json:"code" binding:"required,max=128"`
	Name     string `json:"name" binding:"required,max=64"`
	NodeType int32  `json:"nodeType" binding:"oneof=0 1 2"`
	Action   string `json:"action" binding:"required,max=32"`
	Method   string `json:"method" binding:"omitempty,max=16"`
	Path     string `json:"path" binding:"omitempty,max=255"`
	Sort     int64  `json:"sort"`
	Status   int32  `json:"status" binding:"oneof=0 1"`
	Remark   string `json:"remark" binding:"omitempty,max=500"`
}

// ApiPermissionAssignRequest 授权 API 权限。
type ApiPermissionAssignRequest struct {
	PermissionIds []int64 `json:"permissionIds"`
}
