// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package operlog

import (
	"net/http"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/operlog"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OperLogCleanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogCleanReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := operlog.NewOperLogCleanLogic(r.Context(), svcCtx)
		resp, err := l.OperLogClean(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
