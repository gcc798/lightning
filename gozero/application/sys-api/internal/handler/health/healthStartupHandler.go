// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package health

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/health"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func HealthStartupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewHealthStartupLogic(r.Context(), svcCtx)
		resp, err := l.HealthStartup()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
