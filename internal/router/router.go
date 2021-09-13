package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/towelong/lin-cms-go/api/app/cms"
	v1 "github.com/towelong/lin-cms-go/api/app/v1"
	"github.com/towelong/lin-cms-go/pkg/response"
	"net/http"
)

var _ IRouter = (*Router)(nil)

var Set = wire.NewSet(
	wire.Struct(new(Router), "*"),
	wire.Bind(new(IRouter), new(*Router)),
)

type IRouter interface {
	RegisterAPI(app *gin.Engine)
}

type Router struct {
	AdminAPI *cms.AdminAPI
	UserAPI  *cms.UserAPI
	BookAPI  *v1.BookAPI
}

func (r *Router) RegisterAPI(app *gin.Engine) {

	app.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, response.UnifyResponse(10080, ctx))
	})

	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, response.UnifyResponse(10025, ctx))
	})

	cmsGroup := app.Group("/cms")
	{
		r.UserAPI.RegisterServer(cmsGroup)
		r.AdminAPI.RegisterServer(cmsGroup)
	}

	v1Group := app.Group("/v1")
	{
		r.BookAPI.RegisterServer(v1Group)
	}
}
