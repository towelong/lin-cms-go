package router

import (
	"github.com/gin-gonic/gin"
)

var RouteMetaInfo = make(map[string]Meta)

type Meta struct {
	Module     string `json:"module"`
	Permission string `json:"permission"`
	Mount      bool   `json:"mount"`
}

type LinRouter struct {
	prefix string
	module string
	router *gin.RouterGroup
}

func NewLinRouter(prefix, module string, router *gin.RouterGroup) *LinRouter {
	return &LinRouter{
		prefix: prefix,
		module: module,
		router: router,
	}
}

func (l *LinRouter) Permission(permission string, mount bool) Meta {
	return Meta{
		Permission: permission,
		Module:     l.module,
		Mount:      mount,
	}
}

func (l *LinRouter) LinGET(name string, path string, meta Meta, handlers ...gin.HandlerFunc) gin.IRoutes {
	endpoint := "GET " + name
	RouteMetaInfo[endpoint] = meta
	newHandles := make([]gin.HandlerFunc, 1)
	newHandles[0] = func(ctx *gin.Context) {
		ctx.Set("meta", meta)
		ctx.Next()
	}
	newHandles = append(newHandles, handlers...)
	return l.router.GET(l.prefix+path, newHandles...)
}

func (l *LinRouter) LinPOST(name string, path string, meta Meta, handlers ...gin.HandlerFunc) gin.IRoutes {
	endpoint := "POST " + name
	RouteMetaInfo[endpoint] = meta
	newHandles := make([]gin.HandlerFunc, 1)
	newHandles[0] = func(ctx *gin.Context) {
		ctx.Set("meta", meta)
		ctx.Next()
	}
	newHandles = append(newHandles, handlers...)
	return l.router.POST(l.prefix+path, newHandles...)
}

func (l *LinRouter) LinPUT(name string, path string, meta Meta, handlers ...gin.HandlerFunc) gin.IRoutes {
	endpoint := "PUT " + name
	RouteMetaInfo[endpoint] = meta
	newHandles := make([]gin.HandlerFunc, 1)
	newHandles[0] = func(ctx *gin.Context) {
		ctx.Set("meta", meta)
		ctx.Next()
	}
	newHandles = append(newHandles, handlers...)
	return l.router.PUT(l.prefix+path, newHandles...)
}

func (l *LinRouter) LinDELETE(name string, path string, meta Meta, handlers ...gin.HandlerFunc) gin.IRoutes {
	endpoint := "DELETE " + name
	RouteMetaInfo[endpoint] = meta
	newHandles := make([]gin.HandlerFunc, 1)
	newHandles[0] = func(ctx *gin.Context) {
		ctx.Set("meta", meta)
		ctx.Next()
	}
	newHandles = append(newHandles, handlers...)
	return l.router.DELETE(l.prefix+path, newHandles...)
}

// default http request

func (l *LinRouter) GET(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return l.router.GET(l.prefix+path, handlers...)
}

func (l *LinRouter) POST(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return l.router.POST(l.prefix+path, handlers...)
}

func (l *LinRouter) PUT(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return l.router.PUT(l.prefix+path, handlers...)
}

func (l *LinRouter) DELETE(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return l.router.DELETE(l.prefix+path, handlers...)
}
