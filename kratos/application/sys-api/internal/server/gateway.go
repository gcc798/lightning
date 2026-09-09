package server

import (
	"context"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
)

type AuthGateway interface {
	ValidateAccessToken(ctx context.Context, token string) (*v1.ValidateAccessTokenReply, error)
	CheckPermission(ctx context.Context, userID int64, resource, action string) (*v1.CheckPermissionReply, error)
}

type AttachmentGateway interface {
	Upload(ctx context.Context, req *v1.UploadFileRequest) (*v1.AttachmentItem, error)
	Download(ctx context.Context, id int64) (*v1.AttachmentDownloadReply, error)
}

type GatewayDeps struct {
	Auth       AuthGateway
	Attachment AttachmentGateway
	OperLog    OperLogGateway
	WebSocket  *WebSocketHub
}

type OperLogGateway interface {
	Create(ctx context.Context, item *v1.CreateOperLogRequest) (*v1.MessageReply, error)
}
