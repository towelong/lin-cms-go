package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/router"
)

func InitEngine(r router.IRouter) *gin.Engine {
	app := gin.Default()
	app.Use(middleware.CORS)
	app.Use(middleware.ErrorHandler)
	r.RegisterAPI(app)
	return app
}
