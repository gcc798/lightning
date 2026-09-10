package server

import v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"

const (
	operationHealthReady   = "/custom.health/Ready"
	operationHealthLive    = "/custom.health/Live"
	operationHealthStartup = "/custom.health/Startup"
	operationMetrics       = "/custom.metrics/Expose"
	operationSwagger       = "/custom.swagger/Expose"
)

var publicOperations = map[string]struct{}{
	v1.OperationHealthServicePing:                  {},
	operationHealthReady:                           {},
	operationHealthLive:                            {},
	operationHealthStartup:                         {},
	operationMetrics:                               {},
	v1.OperationAuthServiceLogin:                   {},
	v1.OperationAuthServiceLogout:                  {},
	v1.OperationAuthServiceRefreshToken:            {},
	v1.OperationCaptchaServiceGenerateImageCaptcha: {},
	v1.OperationCaptchaServiceSendSMSCaptcha:       {},
	v1.OperationCaptchaServiceSendEmailCaptcha:     {},
	v1.OperationCaptchaServiceGetEnabledTypes:      {},
}

var operationPermissions = map[string]string{
	v1.OperationMenuServiceGetMenuTree:                     "menu.read",
	v1.OperationMenuServiceGetMenuList:                     "menu.read",
	v1.OperationMenuServiceGetMenuById:                     "menu.read",
	v1.OperationMenuServiceCreateMenu:                      "menu.create",
	v1.OperationMenuServiceUpdateMenu:                      "menu.update",
	v1.OperationMenuServiceDeleteMenu:                      "menu.delete",
	v1.OperationUserServiceCreateUser:                      "user.create",
	v1.OperationUserServiceImportUsers:                     "user.create",
	v1.OperationUserServicePageUser:                        "user.read",
	v1.OperationUserServiceBatchDeleteUser:                 "user.delete",
	v1.OperationUserServiceUpdateUser:                      "user.update",
	v1.OperationUserServiceResetPassword:                   "user.update",
	v1.OperationUserServiceGetUserById:                     "user.read",
	v1.OperationUserServiceDeleteUser:                      "user.delete",
	v1.OperationRoleServiceCreateRole:                      "role.create",
	v1.OperationRoleServiceUpdateRole:                      "role.update",
	v1.OperationRoleServiceDeleteRole:                      "role.delete",
	v1.OperationRoleServiceGetRole:                         "role.read",
	v1.OperationRoleServicePageRole:                        "role.read",
	v1.OperationRoleServiceAssignRoleToUser:                "role.assign",
	v1.OperationRoleServiceRemoveRoleFromUser:              "role.assign",
	v1.OperationRoleServiceGetUserRoles:                    "role.read",
	v1.OperationRoleServiceAddRolePermission:               "role.permission",
	v1.OperationRoleServiceDeleteRolePermission:            "role.permission",
	v1.OperationRoleServiceGetRolePermissions:              "role.permission",
	v1.OperationOrgServiceCreateOrg:                        "org.create",
	v1.OperationOrgServicePageOrg:                          "org.read",
	v1.OperationOrgServiceGetOrgTree:                       "org.read",
	v1.OperationOrgServiceGetOrgById:                       "org.read",
	v1.OperationOrgServiceUpdateOrg:                        "org.update",
	v1.OperationOrgServiceDeleteOrg:                        "org.delete",
	v1.OperationOrgServiceBatchDeleteOrg:                   "org.delete",
	v1.OperationConfigServiceCreateConfig:                  "config.create",
	v1.OperationConfigServicePageConfig:                    "config.read",
	v1.OperationConfigServiceBatchDeleteConfig:             "config.delete",
	v1.OperationConfigServiceGetConfigByCode:               "config.read",
	v1.OperationConfigServiceGetConfigDataByCode:           "config.read",
	v1.OperationConfigServiceUpdateConfig:                  "config.update",
	v1.OperationConfigServiceGetConfigById:                 "config.read",
	v1.OperationConfigServiceDeleteConfig:                  "config.delete",
	v1.OperationDictServiceCreateDict:                      "dict.create",
	v1.OperationDictServicePageDict:                        "dict.read",
	v1.OperationDictServiceBatchDeleteDict:                 "dict.delete",
	v1.OperationDictServiceGetDictByType:                   "dict.read",
	v1.OperationDictServiceGetDictLabel:                    "dict.read",
	v1.OperationDictServiceUpdateDict:                      "dict.update",
	v1.OperationDictServiceGetDictById:                     "dict.read",
	v1.OperationDictServiceDeleteDict:                      "dict.delete",
	v1.OperationLoginLogServiceCreateLoginLog:              "login_log.create",
	v1.OperationLoginLogServicePageLoginLog:                "login_log.read",
	v1.OperationLoginLogServiceBatchDeleteLoginLog:         "login_log.delete",
	v1.OperationLoginLogServiceCleanLoginLog:               "login_log.delete",
	v1.OperationLoginLogServiceUpdateLoginLog:              "login_log.update",
	v1.OperationLoginLogServiceGetLoginLogById:             "login_log.read",
	v1.OperationLoginLogServiceDeleteLoginLog:              "login_log.delete",
	v1.OperationOperLogServiceCreateOperLog:                "oper_log.create",
	v1.OperationOperLogServicePageOperLog:                  "oper_log.read",
	v1.OperationOperLogServiceBatchDeleteOperLog:           "oper_log.delete",
	v1.OperationOperLogServiceCleanOperLog:                 "oper_log.delete",
	v1.OperationOperLogServiceUpdateOperLog:                "oper_log.update",
	v1.OperationOperLogServiceGetOperLogById:               "oper_log.read",
	v1.OperationOperLogServiceDeleteOperLog:                "oper_log.delete",
	v1.OperationStorageEnvServiceCreateStorageEnv:          "storage_env.create",
	v1.OperationStorageEnvServicePageStorageEnv:            "storage_env.read",
	v1.OperationStorageEnvServiceGetDefaultStorageEnv:      "storage_env.read",
	v1.OperationStorageEnvServiceSetDefaultStorageEnv:      "storage_env.manage",
	v1.OperationStorageEnvServiceUpdateStorageEnv:          "storage_env.update",
	v1.OperationStorageEnvServiceGetStorageEnv:             "storage_env.read",
	v1.OperationStorageEnvServiceTestStorageEnvConnection:  "storage_env.read",
	v1.OperationStorageEnvServiceDeleteStorageEnv:          "storage_env.delete",
	v1.OperationAttachmentServiceUploadFile:                "attachment.upload",
	v1.OperationAttachmentServiceBindAttachmentToBusiness:  "attachment.bind",
	v1.OperationAttachmentServiceGetAttachment:             "attachment.read",
	v1.OperationAttachmentServiceListAttachmentsByBusiness: "attachment.read",
	v1.OperationAttachmentServicePageAttachments:           "attachment.read",
	v1.OperationAttachmentServiceDownloadAttachment:        "attachment.download",
	v1.OperationAttachmentServiceGetAttachmentURL:          "attachment.read",
	v1.OperationAttachmentServiceDeleteAttachment:          "attachment.delete",
}

func isPublicOperation(operation string) bool {
	_, ok := publicOperations[operation]
	return ok
}

func permissionForOperation(operation string) (string, bool) {
	permission, ok := operationPermissions[operation]
	return permission, ok
}
