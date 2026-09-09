// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/dict"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DictLabelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DictLabelQueryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dict.NewDictLabelLogic(r.Context(), svcCtx)
		resp, err := l.DictLabel(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
