// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package health

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/health"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func HealthLiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewHealthLiveLogic(r.Context(), svcCtx)
		resp, err := l.HealthLive()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
