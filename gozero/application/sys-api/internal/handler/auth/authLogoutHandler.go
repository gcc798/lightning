package auth

import (
	"net/http"
	"strings"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/auth"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get(svcCtx.Config.Auth.TokenHeader), "Bearer "))
		l := auth.NewAuthLogoutLogic(r.Context(), svcCtx)
		resp, err := l.AuthLogout(token)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
