package middleware

import (
	"net/http"
	"path"
	"strings"

	"github.com/gcc798/microservice-kit/common/auth"
)

type JWTAuthConfig struct {
	Secret      string
	TokenHeader string
	WhiteList   []string
}

type JWTAuthMiddleware struct {
	cfg JWTAuthConfig
}

func NewJWTAuthMiddleware(cfg JWTAuthConfig) *JWTAuthMiddleware {
	if cfg.TokenHeader == "" {
		cfg.TokenHeader = "Authorization"
	}
	return &JWTAuthMiddleware{cfg: cfg}
}

func (m *JWTAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.skip(r.URL.Path) {
			next(w, r)
			return
		}

		token := auth.TokenFromRequest(r, m.cfg.TokenHeader)
		claims, err := auth.ParseAccessToken(token, m.cfg.Secret)
		if err != nil {
			writeJSON(w, CodeUnauthorized, "Token验证失败: "+err.Error(), nil)
			return
		}

		headerClientID := strings.TrimSpace(r.Header.Get("clientId"))
		if headerClientID == "" {
			headerClientID = strings.TrimSpace(r.Header.Get("clientid"))
		}
		if headerClientID == "" {
			headerClientID = strings.TrimSpace(r.URL.Query().Get("clientId"))
		}
		if headerClientID == "" {
			headerClientID = strings.TrimSpace(r.URL.Query().Get("clientid"))
		}
		if headerClientID != "" && claims.ClientID != "" && headerClientID != claims.ClientID {
			writeJSON(w, CodeUnauthorized, "客户端ID与Token不匹配", nil)
			return
		}

		ctx := auth.WithUserContext(r.Context(), auth.UserContext{
			UserID:      claims.UserID,
			UserName:    claims.UserName,
			ClientID:    claims.ClientID,
			DeviceType:  claims.DeviceType,
			OrgID:       claims.OrgID,
			Roles:       claims.Roles,
			Permissions: claims.Permissions,
		})
		next(w, r.WithContext(ctx))
	}
}

func (m *JWTAuthMiddleware) skip(requestPath string) bool {
	cleanPath := path.Clean("/" + strings.TrimSpace(requestPath))
	for _, pattern := range m.cfg.WhiteList {
		if matchPath(pattern, cleanPath) {
			return true
		}
	}
	return false
}

func matchPath(pattern, requestPath string) bool {
	cleanPattern := path.Clean("/" + strings.TrimSpace(pattern))
	if cleanPattern == requestPath {
		return true
	}
	if strings.HasSuffix(cleanPattern, "/*") {
		prefix := strings.TrimSuffix(cleanPattern, "/*")
		return requestPath == prefix || strings.HasPrefix(requestPath, prefix+"/")
	}
	return false
}
