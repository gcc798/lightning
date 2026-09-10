// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"net/http"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/menu"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func MenuListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewMenuListLogic(r.Context(), svcCtx)
		resp, err := l.MenuList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
