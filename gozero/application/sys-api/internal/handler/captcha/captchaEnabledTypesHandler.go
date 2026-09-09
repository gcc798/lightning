// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package captcha

import (
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/captcha"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CaptchaEnabledTypesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := captcha.NewCaptchaEnabledTypesLogic(r.Context(), svcCtx)
		resp, err := l.CaptchaEnabledTypes()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
