// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"net/http"
	"strconv"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/attachment"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AttachmentDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AttachmentIdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := attachment.NewAttachmentDownloadLogic(r.Context(), svcCtx)
		data, err := l.AttachmentDownload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Header().Set("Content-Type", data.ContentType)
			w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(data.FileName))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data.Content)
		}
	}
}
