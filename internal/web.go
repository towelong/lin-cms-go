package internal

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/router"
	"github.com/towelong/lin-cms-go/pkg"
)

func InitEngine(r router.IRouter) *gin.Engine {
	if viper.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.New()
	app.Static("/assets", path.Join(pkg.GetCurrentAbPath(), "/assets"))
	app.Use(gin.Recovery())
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.CORS)
	app.Use(middleware.Log)
	r.RegisterAPI(app)
	return app
}
