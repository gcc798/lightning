// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package storageenv

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/storageenv"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func StorageEnvUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StorageEnvUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := storageenv.NewStorageEnvUpdateLogic(r.Context(), svcCtx)
		resp, err := l.StorageEnvUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
