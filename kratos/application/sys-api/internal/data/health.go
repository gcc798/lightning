package data

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
)

type HealthRepo struct {
	client v1.HealthServiceClient
}

func NewHealthRepo(clients *RPCClientSet) *HealthRepo {
	return &HealthRepo{client: clients.Health}
}

func (r *HealthRepo) Ping(ctx context.Context, name string) (*v1.PingReply, error) {
	return r.client.Ping(ctx, &v1.PingRequest{Name: name})
}
