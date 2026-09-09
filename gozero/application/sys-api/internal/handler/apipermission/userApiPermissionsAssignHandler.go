// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/apipermission"
	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserApiPermissionsAssignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserApiPermissionsAssignReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		ctx := r.Context()
		if userID, err := commonutil.UserIDFromRequest(svcCtx, r); err == nil {
			ctx = commonutil.WithUserID(ctx, userID)
		}
		l := apipermission.NewUserApiPermissionsAssignLogic(ctx, svcCtx)
		resp, err := l.UserApiPermissionsAssign(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
