package commonutil

import (
	"context"

	"github.com/gcc798/lightning/common/auth"
)

func WithUserID(ctx context.Context, userID int64) context.Context {
	return auth.WithUserID(ctx, userID)
}

func UserIDFromContext(ctx context.Context) int64 {
	return auth.UserIDFromContext(ctx)
}

func OrgIDFromContext(ctx context.Context) int64 {
	return auth.OrgIDFromContext(ctx)
}

func RolesFromContext(ctx context.Context) []string {
	return auth.RolesFromContext(ctx)
}

func PermissionsFromContext(ctx context.Context) []string {
	return auth.PermissionsFromContext(ctx)
}
