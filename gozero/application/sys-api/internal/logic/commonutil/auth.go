package commonutil

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	commonauth "github.com/gcc798/microservice-kit/common/auth"
)

func BearerToken(r *http.Request, header string) string {
	token := r.Header.Get(header)
	return strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
}

func UserIDFromRequest(svcCtx *svc.ServiceContext, r *http.Request) (int64, error) {
	token := commonauth.TokenFromRequest(r, svcCtx.Config.Auth.TokenHeader)
	claims, err := commonauth.ParseAccessToken(token, svcCtx.Config.Jwt.Secret)
	if err != nil {
		return 0, fmt.Errorf("未登录")
	}
	return claims.UserID, nil
}
