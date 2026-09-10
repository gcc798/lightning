package sysv1

import (
	"context"

	"github.com/gcc798/microservice-kit/internal/transport"
)

const ServiceName = "sys"

type API interface {
	RecordLogin(context.Context, *RecordLoginRequest) error
	RecordOperations(context.Context, *RecordOperationsRequest) error
	CleanLogs(context.Context, int32) (*CleanLogsResponse, error)
}

type Remote struct{ pool *transport.ClientPool }

func NewRemote(pool *transport.ClientPool) *Remote { return &Remote{pool: pool} }

func (r *Remote) RecordLogin(ctx context.Context, request *RecordLoginRequest) error {
	conn, err := r.pool.Conn(ctx, ServiceName)
	if err != nil {
		return err
	}
	_, err = NewSystemServiceClient(conn).RecordLogin(ctx, request)
	return err
}

func (r *Remote) RecordOperations(ctx context.Context, request *RecordOperationsRequest) error {
	conn, err := r.pool.Conn(ctx, ServiceName)
	if err != nil {
		return err
	}
	_, err = NewSystemServiceClient(conn).RecordOperations(ctx, request)
	return err
}

func (r *Remote) CleanLogs(ctx context.Context, days int32) (*CleanLogsResponse, error) {
	conn, err := r.pool.Conn(ctx, ServiceName)
	if err != nil {
		return nil, err
	}
	return NewSystemServiceClient(conn).CleanLogs(ctx, &CleanLogsRequest{Days: days})
}
