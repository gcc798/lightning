package httpx

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Router keeps route declarations compact while registering native Echo handlers.
type Router struct {
	echo *echo.Echo
}

func NewRouter(e *echo.Echo) *Router { return &Router{echo: e} }

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) { r.echo.ServeHTTP(w, req) }

func (r *Router) Use(middleware ...echo.MiddlewareFunc) { r.echo.Use(middleware...) }

func (r *Router) Group(prefix string, middleware ...echo.MiddlewareFunc) *Group {
	return &Group{group: r.echo.Group(prefix, middleware...)}
}

func (r *Router) GET(path string, handlers ...any)    { register(r.echo.GET, path, handlers...) }
func (r *Router) POST(path string, handlers ...any)   { register(r.echo.POST, path, handlers...) }
func (r *Router) PUT(path string, handlers ...any)    { register(r.echo.PUT, path, handlers...) }
func (r *Router) DELETE(path string, handlers ...any) { register(r.echo.DELETE, path, handlers...) }
func (r *Router) PATCH(path string, handlers ...any)  { register(r.echo.PATCH, path, handlers...) }

type Group struct {
	group *echo.Group
}

func (g *Group) Use(middleware ...echo.MiddlewareFunc) { g.group.Use(middleware...) }
func (g *Group) Group(prefix string, middleware ...echo.MiddlewareFunc) *Group {
	return &Group{group: g.group.Group(prefix, middleware...)}
}
func (g *Group) GET(path string, handlers ...any)    { register(g.group.GET, path, handlers...) }
func (g *Group) POST(path string, handlers ...any)   { register(g.group.POST, path, handlers...) }
func (g *Group) PUT(path string, handlers ...any)    { register(g.group.PUT, path, handlers...) }
func (g *Group) DELETE(path string, handlers ...any) { register(g.group.DELETE, path, handlers...) }
func (g *Group) PATCH(path string, handlers ...any)  { register(g.group.PATCH, path, handlers...) }

type routeRegistrar func(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo

func register(registrar routeRegistrar, path string, handlers ...any) {
	if len(handlers) == 0 {
		panic("httpx: route handler is required")
	}
	middleware := make([]echo.MiddlewareFunc, 0, len(handlers)-1)
	var final echo.HandlerFunc
	for _, candidate := range handlers {
		switch handler := candidate.(type) {
		case echo.MiddlewareFunc:
			middleware = append(middleware, handler)
		case func(echo.HandlerFunc) echo.HandlerFunc:
			middleware = append(middleware, echo.MiddlewareFunc(handler))
		case echo.HandlerFunc:
			if final != nil {
				panic(fmt.Sprintf("httpx: multiple route handlers for %s", path))
			}
			final = handler
		case func(*echo.Context):
			if final != nil {
				panic(fmt.Sprintf("httpx: multiple route handlers for %s", path))
			}
			final = Handler(handler)
		default:
			panic(fmt.Sprintf("httpx: unsupported handler %T for %s", candidate, path))
		}
	}
	if final == nil {
		panic(fmt.Sprintf("httpx: final route handler is required for %s", path))
	}
	registrar(path, final, middleware...)
}
