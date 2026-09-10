// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"net/http"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/auth"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get(svcCtx.Config.Auth.TokenHeader), "Bearer "))
		l := auth.NewLogoutLogic(r.Context(), svcCtx)
		resp, err := l.Logout(token)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
