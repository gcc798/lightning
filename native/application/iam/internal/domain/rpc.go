package iam

import (
	"context"

	iamv1 "github.com/gcc798/microservice-kit/internal/api/iam/v1"
	"github.com/gcc798/microservice-kit/internal/platform/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthzAPI struct {
	tokens      TokenManager
	permissions PermissionService
}

func NewAuthzAPI(tokens TokenManager, permissions PermissionService) *AuthzAPI {
	return &AuthzAPI{tokens: tokens, permissions: permissions}
}

func (a *AuthzAPI) ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error) {
	return a.tokens.ValidateAccessToken(ctx, token)
}

func (a *AuthzAPI) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	return a.permissions.CheckPermission(ctx, userID, resource, action)
}

type GRPCServer struct {
	iamv1.UnimplementedIAMServiceServer
	api iamv1.API
}

func NewGRPCServer(api iamv1.API) *GRPCServer { return &GRPCServer{api: api} }

func (s *GRPCServer) ValidateAccessToken(ctx context.Context, request *iamv1.ValidateAccessTokenRequest) (*iamv1.ValidateAccessTokenResponse, error) {
	claims, err := s.api.ValidateAccessToken(ctx, request.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &iamv1.ValidateAccessTokenResponse{
		UserId: claims.UserId, UserName: claims.UserName, ClientId: claims.ClientId, DeviceType: claims.DeviceType,
	}, nil
}

func (s *GRPCServer) CheckPermission(ctx context.Context, request *iamv1.CheckPermissionRequest) (*iamv1.CheckPermissionResponse, error) {
	allowed, err := s.api.CheckPermission(ctx, request.UserId, request.Resource, request.Action)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &iamv1.CheckPermissionResponse{Allowed: allowed}, nil
}
