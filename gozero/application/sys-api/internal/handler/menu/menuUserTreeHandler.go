// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"net/http"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/menu"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	commonauth "github.com/gcc798/microservice-kit/common/auth"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func MenuUserTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := commonauth.UserIDFromContext(r.Context())

		l := menu.NewMenuUserTreeLogic(r.Context(), svcCtx)
		resp, err := l.MenuUserTree(userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
