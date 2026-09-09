package validator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

func newTestEcho() *echo.Echo {
	e := echo.New()
	e.Validator = EchoValidator{}
	e.Binder = &EchoBinder{}
	Init()
	return e
}

func TestEchoBinderBindsAndValidatesJSON(t *testing.T) {
	e := newTestEcho()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"invalid"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, httptest.NewRecorder())
	var payload struct {
		Email string `json:"email" binding:"required,email"`
	}

	err := c.Bind(&payload)
	require.Error(t, err)
	require.Contains(t, Translate(err), "邮箱格式不正确")
}

func TestEchoBinderBindsQueryTags(t *testing.T) {
	e := newTestEcho()
	req := httptest.NewRequest(http.MethodGet, "/config?code=WechatXcxCfg", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	var payload struct {
		Code string `query:"code" binding:"required"`
	}

	require.NoError(t, c.Bind(&payload))
	require.Equal(t, "WechatXcxCfg", payload.Code)
}
