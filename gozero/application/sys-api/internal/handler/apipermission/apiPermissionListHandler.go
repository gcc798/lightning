// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/apipermission"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ApiPermissionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := apipermission.NewApiPermissionListLogic(r.Context(), svcCtx)
		resp, err := l.ApiPermissionList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
