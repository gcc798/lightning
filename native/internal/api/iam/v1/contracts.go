package iamv1

import (
	"context"

	"github.com/gcc798/lightning/internal/platform/jwt"
	"github.com/gcc798/lightning/internal/transport"
)

const ServiceName = "iam"

type API interface {
	ValidateAccessToken(context.Context, string) (*jwt.Claims, error)
	CheckPermission(context.Context, int64, string, string) (bool, error)
}

type Remote struct{ pool *transport.ClientPool }

func NewRemote(pool *transport.ClientPool) *Remote { return &Remote{pool: pool} }

func (r *Remote) ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error) {
	conn, err := r.pool.Conn(ctx, ServiceName)
	if err != nil {
		return nil, err
	}
	result, err := NewIAMServiceClient(conn).ValidateAccessToken(ctx, &ValidateAccessTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	return &jwt.Claims{UserId: result.UserId, UserName: result.UserName, ClientId: result.ClientId, DeviceType: result.DeviceType}, nil
}

func (r *Remote) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	conn, err := r.pool.Conn(ctx, ServiceName)
	if err != nil {
		return false, err
	}
	result, err := NewIAMServiceClient(conn).CheckPermission(ctx, &CheckPermissionRequest{UserId: userID, Resource: resource, Action: action})
	if err != nil {
		return false, err
	}
	return result.Allowed, nil
}
