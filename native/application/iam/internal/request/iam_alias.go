package request

import "github.com/gcc798/microservice-kit/application/iam/internal/domain"

type ApiPermissionSaveRequest = iam.ApiPermissionSaveRequest
type ApiPermissionAssignRequest = iam.ApiPermissionAssignRequest
type CreateUserRequest = iam.CreateUserRequest
type UpdateUserRequest = iam.UpdateUserRequest
type BatchDeleteUsersRequest = iam.BatchDeleteUsersRequest
type PageUsersRequest = iam.PageUsersRequest
type BatchImportUsersRequest = iam.BatchImportUsersRequest
type ResetPasswordRequest = iam.ResetPasswordRequest
type ChangePasswordRequest = iam.ChangePasswordRequest
type CreateOrgRequest = iam.CreateOrgRequest
type UpdateOrgRequest = iam.UpdateOrgRequest
type BatchDeleteOrgsRequest = iam.BatchDeleteOrgsRequest
type PageOrgsRequest = iam.PageOrgsRequest
type LoginRequest = iam.LoginRequest
type SendSmsCodeRequest = iam.SendSmsCodeRequest
type SendEmailCodeRequest = iam.SendEmailCodeRequest
