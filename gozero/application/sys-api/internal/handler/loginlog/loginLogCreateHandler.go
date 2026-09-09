// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package loginlog

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/loginlog"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginLogCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginLogReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := loginlog.NewLoginLogCreateLogic(r.Context(), svcCtx)
		resp, err := l.LoginLogCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
