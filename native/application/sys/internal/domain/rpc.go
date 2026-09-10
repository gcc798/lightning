package sys

import (
	"context"
	"time"

	"github.com/gcc798/microservice-kit/application/sys/internal/domain/model"
	sysv1 "github.com/gcc798/microservice-kit/internal/api/sys/v1"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type API struct {
	db     *gorm.DB
	logger logging.Logger
}

func NewAPI(db *gorm.DB, logger logging.Logger) *API { return &API{db: db, logger: logger} }

func (a *API) RecordLogin(ctx context.Context, request *sysv1.RecordLoginRequest) error {
	entry := &model.LoginLog{
		ID: request.Id, UserName: request.UserName, Ipaddr: request.IpAddress,
		Browser: request.Browser, Os: request.Os, Status: request.Status, Msg: request.Message,
		LoginTime: utils.LocalTime(time.UnixMilli(request.LoginTimeUnixMilli)), ClientId: request.ClientId,
	}
	return a.db.WithContext(ctx).Create(entry).Error
}

func (a *API) RecordOperations(ctx context.Context, request *sysv1.RecordOperationsRequest) error {
	if len(request.Logs) == 0 {
		return nil
	}
	entries := make([]model.OperLog, 0, len(request.Logs))
	for _, item := range request.Logs {
		entries = append(entries, model.OperLog{
			ID: item.Id, Title: item.Title, BusinessType: item.BusinessType, Method: item.Method,
			RequestMethod: item.RequestMethod, DeviceType: item.DeviceType, OperName: item.OperatorName,
			OperUrl: item.Url, OperIp: item.IpAddress, OperParam: item.Parameters, Status: item.Status,
			ErrorMsg: item.ErrorMessage, OperTime: utils.LocalTime(time.UnixMilli(item.OperationTimeUnixMilli)),
			CostTime: item.CostMillis, UserAgent: item.UserAgent,
		})
	}
	return a.db.WithContext(ctx).Create(&entries).Error
}

func (a *API) CleanLogs(ctx context.Context, days int32) (*sysv1.CleanLogsResponse, error) {
	if days <= 0 {
		return nil, status.Error(codes.InvalidArgument, "days must be positive")
	}
	cutoff := time.Now().AddDate(0, 0, -int(days))
	login := a.db.WithContext(ctx).Where("login_time < ?", cutoff).Delete(&model.LoginLog{})
	if login.Error != nil {
		return nil, login.Error
	}
	operation := a.db.WithContext(ctx).Where("oper_time < ?", cutoff).Delete(&model.OperLog{})
	if operation.Error != nil {
		return nil, operation.Error
	}
	a.logger.Info("cleaned system logs", zap.Int64("loginLogs", login.RowsAffected), zap.Int64("operationLogs", operation.RowsAffected))
	return &sysv1.CleanLogsResponse{LoginLogs: login.RowsAffected, OperationLogs: operation.RowsAffected}, nil
}

type GRPCServer struct {
	sysv1.UnimplementedSystemServiceServer
	api sysv1.API
}

func NewGRPCServer(api sysv1.API) *GRPCServer { return &GRPCServer{api: api} }

func (s *GRPCServer) RecordLogin(ctx context.Context, request *sysv1.RecordLoginRequest) (*sysv1.Empty, error) {
	if err := s.api.RecordLogin(ctx, request); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sysv1.Empty{}, nil
}

func (s *GRPCServer) RecordOperations(ctx context.Context, request *sysv1.RecordOperationsRequest) (*sysv1.Empty, error) {
	if err := s.api.RecordOperations(ctx, request); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sysv1.Empty{}, nil
}

func (s *GRPCServer) CleanLogs(ctx context.Context, request *sysv1.CleanLogsRequest) (*sysv1.CleanLogsResponse, error) {
	return s.api.CleanLogs(ctx, request.Days)
}
