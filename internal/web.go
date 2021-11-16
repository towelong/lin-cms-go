package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/router"
)

func InitEngine(r router.IRouter) *gin.Engine {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.CORS)
	app.Use(middleware.Logger)
	r.RegisterAPI(app)
	return app
}
