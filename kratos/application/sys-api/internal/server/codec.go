package server

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/protobuf/types/known/structpb"
)

const successCode = 200

type successEnvelope struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type pageEnvelope struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Rows  any    `json:"rows"`
	Total int64  `json:"total"`
}

type successProfile struct {
	msg          string
	data         any
	overrideData bool
}

var operationSuccessProfiles = map[string]successProfile{
	v1.OperationAuthServiceLogout:                          {msg: "success", data: "ok", overrideData: true},
	v1.OperationUserServiceCreateUser:                      {msg: "success", data: map[string]any{"userId": "ok"}, overrideData: true},
	v1.OperationUserServiceUpdateUser:                      {msg: "success", data: "ok", overrideData: true},
	v1.OperationUserServiceDeleteUser:                      {msg: "success", data: "ok", overrideData: true},
	v1.OperationUserServiceBatchDeleteUser:                 {msg: "success", data: "ok", overrideData: true},
	v1.OperationUserServiceResetPassword:                   {msg: "success", data: "ok", overrideData: true},
	v1.OperationUserServiceChangePassword:                  {msg: "success", data: "ok", overrideData: true},
	v1.OperationOrgServiceUpdateOrg:                        {msg: "success", data: "ok", overrideData: true},
	v1.OperationOrgServiceDeleteOrg:                        {msg: "success", data: "ok", overrideData: true},
	v1.OperationOrgServiceBatchDeleteOrg:                   {msg: "success", data: "ok", overrideData: true},
	v1.OperationConfigServiceCreateConfig:                  {msg: "创建配置成功", data: nil, overrideData: true},
	v1.OperationConfigServiceUpdateConfig:                  {msg: "更新配置成功", data: nil, overrideData: true},
	v1.OperationConfigServiceDeleteConfig:                  {msg: "删除配置成功", data: nil, overrideData: true},
	v1.OperationConfigServiceBatchDeleteConfig:             {msg: "批量删除配置成功", data: nil, overrideData: true},
	v1.OperationDictServiceCreateDict:                      {msg: "创建字典成功", data: nil, overrideData: true},
	v1.OperationDictServiceUpdateDict:                      {msg: "更新字典成功", data: nil, overrideData: true},
	v1.OperationDictServiceDeleteDict:                      {msg: "删除字典成功", data: nil, overrideData: true},
	v1.OperationDictServiceBatchDeleteDict:                 {msg: "批量删除字典成功", data: nil, overrideData: true},
	v1.OperationLoginLogServiceCreateLoginLog:              {msg: "创建登录日志成功", data: nil, overrideData: true},
	v1.OperationLoginLogServiceUpdateLoginLog:              {msg: "更新登录日志成功", data: nil, overrideData: true},
	v1.OperationLoginLogServiceDeleteLoginLog:              {msg: "删除登录日志成功", data: nil, overrideData: true},
	v1.OperationLoginLogServiceBatchDeleteLoginLog:         {msg: "批量删除登录日志成功", data: nil, overrideData: true},
	v1.OperationOperLogServiceCreateOperLog:                {msg: "创建操作日志成功", data: nil, overrideData: true},
	v1.OperationOperLogServiceUpdateOperLog:                {msg: "更新操作日志成功", data: nil, overrideData: true},
	v1.OperationOperLogServiceDeleteOperLog:                {msg: "删除操作日志成功", data: nil, overrideData: true},
	v1.OperationOperLogServiceBatchDeleteOperLog:           {msg: "批量删除操作日志成功", data: nil, overrideData: true},
	v1.OperationStorageEnvServiceCreateStorageEnv:          {msg: "创建存储环境成功"},
	v1.OperationStorageEnvServiceUpdateStorageEnv:          {msg: "更新存储环境成功", data: nil, overrideData: true},
	v1.OperationStorageEnvServiceDeleteStorageEnv:          {msg: "删除存储环境成功", data: nil, overrideData: true},
	v1.OperationStorageEnvServiceSetDefaultStorageEnv:      {msg: "设置默认存储环境成功", data: nil, overrideData: true},
	v1.OperationStorageEnvServiceGetStorageEnv:             {msg: "获取存储环境详情成功"},
	v1.OperationStorageEnvServiceGetDefaultStorageEnv:      {msg: "获取默认存储环境成功"},
	v1.OperationStorageEnvServiceTestStorageEnvConnection:  {msg: "存储环境连接测试成功", data: nil, overrideData: true},
	v1.OperationAttachmentServiceUploadFile:                {msg: "上传文件成功"},
	v1.OperationAttachmentServiceBindAttachmentToBusiness:  {msg: "绑定附件到业务成功", data: nil, overrideData: true},
	v1.OperationAttachmentServiceGetAttachment:             {msg: "获取附件详情成功"},
	v1.OperationAttachmentServiceListAttachmentsByBusiness: {msg: "查询业务附件列表成功"},
	v1.OperationAttachmentServiceGetAttachmentURL:          {msg: "获取附件URL成功"},
	v1.OperationAttachmentServiceDeleteAttachment:          {msg: "删除附件成功", data: nil, overrideData: true},
}

func responseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	msg, data, isPage, total := buildSuccessPayload(r, v)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if isPage {
		return encoder.Encode(pageEnvelope{Code: successCode, Msg: msg, Rows: data, Total: total})
	}
	return encoder.Encode(successEnvelope{Code: successCode, Msg: msg, Data: data})
}

func errorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	se := kerrors.FromError(err)
	code := int(se.Code)
	if code == 0 {
		code = http.StatusInternalServerError
	}
	payload := map[string]any{"code": code, "msg": normalizeMessage(se.Message), "data": nil}
	_ = json.NewEncoder(w).Encode(payload)
}

func buildSuccessPayload(r *http.Request, v any) (string, any, bool, int64) {
	if v == nil {
		msg, data := applySuccessProfile(r, "success", nil)
		return msg, data, false, 0
	}
	if pageRows, total, ok := extractPagePayload(r, v); ok {
		msg, data := applySuccessProfile(r, "success", pageRows)
		return msg, data, true, total
	}
	switch reply := v.(type) {
	case *v1.MessageReply:
		defaultMsg := "success"
		defaultData := any(nil)
		if profile, ok := successProfileForRequest(r); ok {
			if profile.msg != "" {
				defaultMsg = profile.msg
			}
			if profile.overrideData {
				defaultData = profile.data
			}
			return defaultMsg, defaultData, false, 0
		}
		message := strings.TrimSpace(reply.GetMessage())
		if strings.EqualFold(message, "ok") || message == "" {
			return defaultMsg, defaultData, false, 0
		}
		return normalizeMessage(message), defaultData, false, 0
	case *v1.LoginReply:
		msg, data := applySuccessProfile(r, "success", adaptLoginReply(reply))
		return msg, data, false, 0
	case *v1.RefreshTokenReply:
		msg, data := applySuccessProfile(r, "success", adaptRefreshReply(reply))
		return msg, data, false, 0
	case *v1.CaptchaReply:
		msg, data := applySuccessProfile(r, "success", adaptCaptchaReply(reply))
		return msg, data, false, 0
	case *v1.LogCleanReply:
		msg, data := applySuccessProfile(r, normalizeMessage(reply.GetMessage()), map[string]any{"count": reply.GetCount(), "days": reply.GetDays()})
		return msg, data, false, 0
	case *v1.GetEnabledTypesReply:
		msg, data := applySuccessProfile(r, "success", reply.GetItems())
		return msg, data, false, 0
	case *v1.AttachmentListReply:
		msg, data := applySuccessProfile(r, "success", adaptValue(r, reply.GetItems()))
		return msg, data, false, 0
	case *v1.MenuTreeReply:
		msg, data := applySuccessProfile(r, "success", adaptValue(r, reply.GetItems()))
		return msg, data, false, 0
	case *v1.MenuListReply:
		msg, data := applySuccessProfile(r, "success", adaptValue(r, reply.GetItems()))
		return msg, data, false, 0
	case *v1.ConfigListReply:
		msg, data := applySuccessProfile(r, "success", adaptValue(r, reply.GetItems()))
		return msg, data, false, 0
	case *v1.ConfigDataReply:
		msg, data := applySuccessProfile(r, "success", map[string]any{"code": configCode(r), "data": protoValueInterface(reply.GetData())})
		return msg, data, false, 0
	case *v1.DictListReply:
		msg, data := applySuccessProfile(r, "success", adaptValue(r, reply.GetItems()))
		return msg, data, false, 0
	case *v1.DictLabelReply:
		msg, data := applySuccessProfile(r, "success", reply.GetLabel())
		return msg, data, false, 0
	case *v1.ImportUsersReply:
		msg, data := applySuccessProfile(r, "success", map[string]any{"successCount": reply.GetSuccessCount(), "failCount": reply.GetFailCount(), "errors": reply.GetErrors()})
		return msg, data, false, 0
	default:
		msg, data := applySuccessProfile(r, "success", adaptValue(r, v))
		return msg, data, false, 0
	}
}

func extractPagePayload(r *http.Request, v any) (any, int64, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, 0, false
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, 0, false
		}
	}
	getList := reflect.ValueOf(v).MethodByName("GetList")
	getTotal := reflect.ValueOf(v).MethodByName("GetTotal")
	if !getList.IsValid() || !getTotal.IsValid() {
		return nil, 0, false
	}
	listValues := getList.Call(nil)
	totalValues := getTotal.Call(nil)
	if len(listValues) != 1 || len(totalValues) != 1 {
		return nil, 0, false
	}
	total, ok := totalValues[0].Interface().(int64)
	if !ok {
		return nil, 0, false
	}
	return adaptValue(r, listValues[0].Interface()), total, true
}

func adaptValue(r *http.Request, v any) any {
	switch item := v.(type) {
	case nil:
		return nil
	case *v1.UserInfo:
		return map[string]any{"userId": item.GetUserId(), "username": item.GetUsername(), "nickname": item.GetNickname(), "phonenumber": item.GetPhonenumber(), "email": item.GetEmail(), "avatar": item.GetAvatar(), "userType": item.GetUserType()}
	case *v1.UserItem:
		return map[string]any{"id": item.GetUserId(), "userId": item.GetUserId(), "userName": item.GetUserName(), "username": item.GetUserName(), "nickName": item.GetNickName(), "nickname": item.GetNickName(), "userType": item.GetUserType(), "email": item.GetEmail(), "phonenumber": item.GetPhonenumber(), "sex": item.GetSex(), "avatar": item.GetAvatar(), "status": item.GetStatus(), "sort": item.GetSort(), "loginIp": item.GetLoginIp(), "loginDate": item.GetLoginDate(), "openId": item.GetOpenId(), "unionId": item.GetUnionId(), "remark": item.GetRemark(), "createBy": item.GetCreateBy(), "updateBy": item.GetUpdateBy(), "createdAt": item.GetCreatedTime(), "createdTime": item.GetCreatedTime(), "updatedAt": item.GetUpdatedTime(), "updatedTime": item.GetUpdatedTime()}
	case *v1.RoleItem:
		return map[string]any{"id": item.GetRoleId(), "roleId": item.GetRoleId(), "roleKey": item.GetRoleKey(), "roleName": item.GetRoleName(), "sort": item.GetSort(), "status": item.GetStatus(), "dataScope": item.GetDataScope(), "isSystem": item.GetIsSystem(), "remark": item.GetRemark(), "createBy": item.GetCreateBy(), "createTime": item.GetCreatedTime(), "createdTime": item.GetCreatedTime(), "updateTime": item.GetUpdatedTime(), "updatedTime": item.GetUpdatedTime()}
	case *v1.MenuItem:
		children := adaptValue(r, item.GetChildren())
		return map[string]any{"id": item.GetId(), "menuName": item.GetMenuName(), "parentId": item.GetParentId(), "sort": item.GetSort(), "path": item.GetPath(), "component": item.GetComponent(), "query": item.GetQuery(), "isFrame": item.GetIsFrame(), "isCache": item.GetIsCache(), "menuType": item.GetMenuType(), "visible": item.GetVisible(), "status": item.GetStatus(), "perms": item.GetPerms(), "icon": item.GetIcon(), "remark": item.GetRemark(), "createBy": item.GetCreateBy(), "updateBy": item.GetUpdateBy(), "createdTime": item.GetCreatedTime(), "updatedTime": item.GetUpdatedTime(), "children": children}
	case *v1.OrgItem:
		children := adaptValue(r, item.GetChildren())
		return map[string]any{"id": item.GetOrgId(), "orgId": item.GetOrgId(), "parentId": item.GetParentId(), "orgName": item.GetOrgName(), "orgCode": item.GetOrgCode(), "orgType": item.GetOrgType(), "leader": item.GetLeader(), "phone": item.GetPhone(), "email": item.GetEmail(), "status": item.GetStatus(), "sort": item.GetSort(), "orderNum": item.GetSort(), "remark": item.GetRemark(), "createBy": item.GetCreateBy(), "updateBy": item.GetUpdateBy(), "createTime": item.GetCreatedTime(), "createdTime": item.GetCreatedTime(), "updateTime": item.GetUpdatedTime(), "updatedTime": item.GetUpdatedTime(), "children": children}
	case *v1.ConfigItem:
		configData := protoValueInterface(item.GetData())
		return map[string]any{"id": item.GetId(), "name": item.GetName(), "code": item.GetCode(), "data": configData, "value": configData, "status": item.GetStatus(), "remark": item.GetRemark(), "createBy": item.GetCreateBy(), "updateBy": item.GetUpdateBy(), "createdTime": item.GetCreatedTime(), "updatedTime": item.GetUpdatedTime()}
	case *v1.DictItem:
		return map[string]any{"id": item.GetId(), "parentId": item.GetParentId(), "dictType": item.GetDictType(), "dictLabel": item.GetDictLabel(), "dictValue": item.GetDictValue(), "isDefault": item.GetIsDefault(), "status": item.GetStatus(), "sort": item.GetSort(), "remark": item.GetRemark(), "createBy": item.GetCreateBy(), "updateBy": item.GetUpdateBy(), "createdTime": item.GetCreatedTime(), "updatedTime": item.GetUpdatedTime()}
	case *v1.StorageEnvItem:
		return map[string]any{"id": item.GetId(), "name": item.GetName(), "code": item.GetCode(), "storageType": item.GetStorageType(), "isDefault": item.GetIsDefault(), "status": item.GetStatus(), "config": protoStructInterface(item.GetConfig()), "remark": item.GetRemark(), "createBy": item.GetCreateBy(), "updateBy": item.GetUpdateBy(), "createTime": item.GetCreatedTime(), "createdAt": item.GetCreatedTime(), "createdTime": item.GetCreatedTime(), "updateTime": item.GetUpdatedTime(), "updatedAt": item.GetUpdatedTime(), "updatedTime": item.GetUpdatedTime()}
	case *v1.AttachmentItem:
		return map[string]any{"id": item.GetAttachmentId(), "attachmentId": item.GetAttachmentId(), "envId": item.GetEnvId(), "storageEnvId": item.GetEnvId(), "fileName": item.GetFileName(), "fileKey": item.GetFileKey(), "filePath": item.GetFileKey(), "fileSize": item.GetFileSize(), "fileType": item.GetFileType(), "fileExt": item.GetFileExt(), "businessType": item.GetBusinessType(), "businessId": item.GetBusinessId(), "businessField": item.GetBusinessField(), "isPublic": item.GetIsPublic(), "accessUrl": item.GetAccessUrl(), "fileUrl": item.GetAccessUrl(), "url": item.GetAccessUrl(), "metadata": item.GetMetadata(), "status": item.GetStatus(), "expireTime": item.GetExpireTime(), "createBy": item.GetCreateBy(), "uploadBy": item.GetCreateBy(), "createTime": item.GetCreateTime(), "uploadTime": item.GetCreateTime(), "updateTime": item.GetUpdateTime(), "updatedTime": item.GetUpdateTime()}
	case *v1.AttachmentURLReply:
		expires := item.GetExpires()
		if expires <= 0 {
			expires = attachmentURLExpires(r)
		}
		return map[string]any{"attachmentId": item.GetAttachmentId(), "url": item.GetUrl(), "expires": expires}
	case *v1.LogItem:
		return adaptLogItem(r, item)
	}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		return adaptValue(r, rv.Elem().Interface())
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		items := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			items = append(items, adaptValue(r, rv.Index(i).Interface()))
		}
		return items
	}
	return v
}

func adaptLoginReply(reply *v1.LoginReply) map[string]any {
	return map[string]any{
		"access_token":       reply.GetAccessToken(),
		"refresh_token":      reply.GetRefreshToken(),
		"expires_in":         reply.GetExpiresIn(),
		"refresh_expires_in": reply.GetRefreshExpiresIn(),
		"user_info":          adaptValue(nil, reply.GetUserInfo()),
	}
}

func adaptRefreshReply(reply *v1.RefreshTokenReply) map[string]any {
	return map[string]any{
		"access_token":       reply.GetAccessToken(),
		"refresh_token":      reply.GetRefreshToken(),
		"expires_in":         reply.GetExpiresIn(),
		"refresh_expires_in": reply.GetRefreshExpiresIn(),
	}
}

func adaptCaptchaReply(reply *v1.CaptchaReply) map[string]any {
	data := map[string]any{}
	if reply.GetData() != nil {
		data = reply.GetData().AsMap()
	}
	return map[string]any{
		"id":       reply.GetId(),
		"type":     reply.GetType(),
		"data":     data,
		"expireAt": reply.GetExpireAt(),
	}
}

func adaptLogItem(r *http.Request, item *v1.LogItem) map[string]any {
	operation := requestOperation(r)
	if strings.Contains(operation, "LoginLog") || strings.Contains(r.URL.Path, "/loginLog") {
		return map[string]any{"id": item.GetId(), "userName": item.GetUserName(), "ipaddr": item.GetIpaddr(), "loginLocation": item.GetLoginLocation(), "browser": item.GetBrowser(), "os": item.GetOs(), "status": parseStatusInt(item.GetStatus()), "msg": item.GetMsg(), "loginTime": item.GetCreatedTime(), "clientId": item.GetClientId()}
	}
	return map[string]any{"id": item.GetId(), "title": item.GetTitle(), "businessType": item.GetBusinessType(), "method": item.GetMethod(), "requestMethod": item.GetRequestMethod(), "operatorType": "", "operName": item.GetUserName(), "operUrl": item.GetOperUrl(), "operIp": item.GetIpaddr(), "operLocation": item.GetOperLocation(), "operParam": item.GetOperParam(), "jsonResult": item.GetJsonResult(), "status": item.GetStatus(), "errorMsg": item.GetMsg(), "operTime": item.GetCreatedTime(), "costTime": item.GetCostTime(), "userAgent": item.GetUserAgent(), "deviceType": item.GetDeviceType()}
}

func requestOperation(r *http.Request) string {
	if r == nil {
		return ""
	}
	if tr, ok := transport.FromServerContext(r.Context()); ok {
		return strings.TrimSpace(tr.Operation())
	}
	return ""
}

func successProfileForRequest(r *http.Request) (successProfile, bool) {
	if r == nil {
		return successProfile{}, false
	}
	profile, ok := operationSuccessProfiles[requestOperation(r)]
	return profile, ok
}

func applySuccessProfile(r *http.Request, msg string, data any) (string, any) {
	profile, ok := successProfileForRequest(r)
	if !ok {
		return msg, data
	}
	if profile.msg != "" {
		msg = profile.msg
	}
	if profile.overrideData {
		data = profile.data
	}
	return msg, data
}

func normalizeMessage(msg string) string {
	msg = strings.TrimSpace(msg)
	if msg == "" || strings.EqualFold(msg, "ok") {
		return "success"
	}
	return msg
}

func parseStatusInt(value string) int {
	value = strings.TrimSpace(value)
	switch value {
	case "1":
		return 1
	default:
		return 0
	}
}

func attachmentURLExpires(r *http.Request) int64 {
	if r == nil {
		return 3600
	}
	value := strings.TrimSpace(r.URL.Query().Get("expires"))
	if value == "" {
		return 3600
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 3600
	}
	return parsed
}

func decodeJSONString(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	var out any
	if err := json.Unmarshal([]byte(value), &out); err != nil {
		return value
	}
	return out
}

func protoValueInterface(value *structpb.Value) any {
	if value == nil {
		return nil
	}
	return value.AsInterface()
}

func protoStructInterface(value *structpb.Struct) any {
	if value == nil {
		return nil
	}
	return value.AsMap()
}

func configCode(r *http.Request) string {
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.URL.Query().Get("code"))
}
