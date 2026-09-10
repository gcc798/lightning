package data

import (
	"context"
	"strings"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
)

type AuthRepo struct {
	client v1.AuthServiceClient
}

func NewAuthRepo(clients *RPCClientSet) *AuthRepo {
	return &AuthRepo{client: clients.Auth}
}

func (r *AuthRepo) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return r.client.Login(ctx, req)
}

func (r *AuthRepo) Logout(ctx context.Context) (*v1.MessageReply, error) {
	return r.client.Logout(ctx, &v1.LogoutRequest{})
}

func (r *AuthRepo) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	return r.client.RefreshToken(ctx, req)
}

func (r *AuthRepo) ValidateAccessToken(ctx context.Context, token string) (*v1.ValidateAccessTokenReply, error) {
	return r.client.ValidateAccessToken(ctx, &v1.ValidateAccessTokenRequest{AccessToken: strings.TrimSpace(token)})
}

func (r *AuthRepo) CheckPermission(ctx context.Context, userID int64, resource, action string) (*v1.CheckPermissionReply, error) {
	return r.client.CheckPermission(ctx, &v1.CheckPermissionRequest{UserId: userID, Resource: resource, Action: action})
}
