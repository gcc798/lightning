package resourcev1

import (
	"context"

	"github.com/gcc798/lightning/internal/transport"
)

const ServiceName = "resource"

type API interface {
	CleanExpired(context.Context) (*CleanExpiredResponse, error)
}

type Remote struct{ pool *transport.ClientPool }

func NewRemote(pool *transport.ClientPool) *Remote { return &Remote{pool: pool} }

func (r *Remote) CleanExpired(ctx context.Context) (*CleanExpiredResponse, error) {
	conn, err := r.pool.Conn(ctx, ServiceName)
	if err != nil {
		return nil, err
	}
	return NewResourceServiceClient(conn).CleanExpired(ctx, &CleanExpiredRequest{})
}
