// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/apipermission"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RoleApiPermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleApiPermissionsPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := apipermission.NewRoleApiPermissionsLogic(r.Context(), svcCtx)
		resp, err := l.RoleApiPermissions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
