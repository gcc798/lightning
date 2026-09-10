package data

import (
	"context"
	"errors"
	"strings"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent/loginlog"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent/operlog"
)

func loginLogEntityToItem(item *entpkg.LoginLog) *v1.LogItem {
	if item == nil {
		return nil
	}
	return &v1.LogItem{
		Id:            item.ID,
		Title:         "login",
		UserName:      item.UserName,
		Ipaddr:        item.Ipaddr,
		Status:        statusString(int64(item.Status)),
		Msg:           item.Msg,
		CreatedTime:   formatTime(item.LoginTime),
		LoginLocation: item.LoginLocation,
		Browser:       item.Browser,
		Os:            item.Os,
		ClientId:      item.ClientID,
	}
}

func operLogEntityToItem(item *entpkg.OperLog) *v1.LogItem {
	if item == nil {
		return nil
	}
	return &v1.LogItem{
		Id:            item.ID,
		Title:         item.Title,
		UserName:      item.OperName,
		Ipaddr:        item.OperIP,
		Status:        item.Status,
		Msg:           item.ErrorMsg,
		CreatedTime:   formatTime(item.OperTime),
		BusinessType:  item.BusinessType,
		Method:        item.Method,
		RequestMethod: item.RequestMethod,
		DeviceType:    item.DeviceType,
		OperUrl:       item.OperURL,
		OperLocation:  item.OperLocation,
		OperParam:     item.OperParam,
		JsonResult:    item.JSONResult,
		CostTime:      item.CostTime,
		UserAgent:     item.UserAgent,
	}
}

func (r *Resources) CreateLoginLog(ctx context.Context, req *v1.CreateLoginLogRequest) error {
	status := req.GetStatus()
	if status != 0 {
		status = 1
	}
	_, err := r.Ent.LoginLog.Create().
		SetID(nextID()).
		SetUserName(req.GetUserName()).
		SetIpaddr(req.GetIpaddr()).
		SetLoginLocation(req.GetLoginLocation()).
		SetBrowser(req.GetBrowser()).
		SetOs(req.GetOs()).
		SetStatus(status).
		SetMsg(req.GetMsg()).
		SetLoginTime(time.Now()).
		SetClientID(req.GetClientId()).
		Save(ctx)
	return err
}

func (r *Resources) CreateOperLog(ctx context.Context, req *v1.CreateOperLogRequest) error {
	_, err := r.Ent.OperLog.Create().
		SetID(nextID()).
		SetTitle(req.GetTitle()).
		SetBusinessType(req.GetBusinessType()).
		SetMethod(req.GetMethod()).
		SetRequestMethod(req.GetRequestMethod()).
		SetDeviceType(req.GetDeviceType()).
		SetOperName(req.GetOperName()).
		SetOperURL(req.GetOperUrl()).
		SetOperIP(req.GetOperIp()).
		SetOperLocation(req.GetOperLocation()).
		SetOperParam(req.GetOperParam()).
		SetJSONResult(req.GetJsonResult()).
		SetStatus(req.GetStatus()).
		SetErrorMsg(req.GetErrorMsg()).
		SetOperTime(time.Now()).
		SetCostTime(req.GetCostTime()).
		SetUserAgent(req.GetUserAgent()).
		Save(ctx)
	return err
}

func (r *Resources) PageLoginLogs(ctx context.Context, req *v1.PageLoginLogRequest) (*v1.PageLogReply, error) {
	items, err := r.Ent.LoginLog.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.LogItem, 0, len(items))
	startTime := parseQueryTime(req.GetStartTime())
	endTime := parseQueryTime(req.GetEndTime())
	for _, item := range items {
		if req.GetUserName() != "" && !strings.Contains(item.UserName, req.GetUserName()) {
			continue
		}
		if req.GetIpaddr() != "" && !strings.Contains(item.Ipaddr, req.GetIpaddr()) {
			continue
		}
		if req.Status != nil && item.Status != req.GetStatus() {
			continue
		}
		if !inTimeRange(item.LoginTime, startTime, endTime) {
			continue
		}
		filtered = append(filtered, loginLogEntityToItem(item))
	}
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageLogReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}

func (r *Resources) PageOperLogs(ctx context.Context, req *v1.PageOperLogRequest) (*v1.PageLogReply, error) {
	items, err := r.Ent.OperLog.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.LogItem, 0, len(items))
	startTime := parseQueryTime(req.GetStartTime())
	endTime := parseQueryTime(req.GetEndTime())
	for _, item := range items {
		if req.GetTitle() != "" && !strings.Contains(item.Title, req.GetTitle()) {
			continue
		}
		if req.GetOperName() != "" && !strings.Contains(item.OperName, req.GetOperName()) {
			continue
		}
		if req.GetBusinessType() != "" && item.BusinessType != req.GetBusinessType() {
			continue
		}
		if req.Status != nil && item.Status != req.GetStatus() {
			continue
		}
		if !inTimeRange(item.OperTime, startTime, endTime) {
			continue
		}
		filtered = append(filtered, operLogEntityToItem(item))
	}
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageLogReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}

func (r *Resources) GetLog(ctx context.Context, kind string, id int64) (*v1.LogItem, error) {
	if kind == "oper" {
		item, err := r.Ent.OperLog.Get(ctx, id)
		if err != nil {
			if entpkg.IsNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return operLogEntityToItem(item), nil
	}
	item, err := r.Ent.LoginLog.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return loginLogEntityToItem(item), nil
}

func (r *Resources) UpdateLoginLog(ctx context.Context, req *v1.UpdateLoginLogRequest) error {
	status := req.GetStatus()
	if status != 0 {
		status = 1
	}
	_, err := r.Ent.LoginLog.UpdateOneID(req.GetId()).
		SetUserName(req.GetUserName()).
		SetIpaddr(req.GetIpaddr()).
		SetLoginLocation(req.GetLoginLocation()).
		SetBrowser(req.GetBrowser()).
		SetOs(req.GetOs()).
		SetStatus(status).
		SetMsg(req.GetMsg()).
		SetClientID(req.GetClientId()).
		Save(ctx)
	return err
}

func (r *Resources) UpdateOperLog(ctx context.Context, req *v1.UpdateOperLogRequest) error {
	_, err := r.Ent.OperLog.UpdateOneID(req.GetId()).
		SetTitle(req.GetTitle()).
		SetBusinessType(req.GetBusinessType()).
		SetMethod(req.GetMethod()).
		SetRequestMethod(req.GetRequestMethod()).
		SetDeviceType(req.GetDeviceType()).
		SetOperName(req.GetOperName()).
		SetOperURL(req.GetOperUrl()).
		SetOperIP(req.GetOperIp()).
		SetOperLocation(req.GetOperLocation()).
		SetOperParam(req.GetOperParam()).
		SetJSONResult(req.GetJsonResult()).
		SetStatus(req.GetStatus()).
		SetErrorMsg(req.GetErrorMsg()).
		SetCostTime(req.GetCostTime()).
		SetUserAgent(req.GetUserAgent()).
		Save(ctx)
	return err
}

func (r *Resources) DeleteLogs(ctx context.Context, kind string, ids ...int64) error {
	for _, id := range ids {
		if kind == "oper" {
			if err := r.Ent.OperLog.DeleteOneID(id).Exec(ctx); err != nil && !entpkg.IsNotFound(err) {
				return err
			}
		} else {
			if err := r.Ent.LoginLog.DeleteOneID(id).Exec(ctx); err != nil && !entpkg.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}

func (r *Resources) CleanLogs(ctx context.Context, kind string, days int) (int64, error) {
	if days <= 0 {
		return 0, errors.New("清理天数必须大于等于1天")
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	if kind == "oper" {
		count, err := r.Ent.OperLog.Delete().
			Where(operlog.OperTimeLT(cutoff)).
			Exec(ctx)
		return int64(count), err
	}
	count, err := r.Ent.LoginLog.Delete().
		Where(loginlog.LoginTimeLT(cutoff)).
		Exec(ctx)
	return int64(count), err
}
