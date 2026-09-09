// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/attachment"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AttachmentBindHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AttachmentBindReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := attachment.NewAttachmentBindLogic(r.Context(), svcCtx)
		resp, err := l.AttachmentBind(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
