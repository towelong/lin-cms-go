package cms

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/extension/file"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
)

type FileAPI struct {
	LocalUploader file.Uploader
	Auth          middleware.Auth
}

func (f FileAPI) UploadFile(ctx *gin.Context) {
	files, err := f.LocalUploader.Upload(ctx)
	if err != nil {
		ctx.Error(response.NewResponse(10210))
		return
	}
	ctx.JSON(http.StatusOK, files)
}

func (f FileAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	fileRouter := router.NewLinRouter("/file", "文件", routerGroup)
	fileRouter.POST("", f.Auth.LoginRequired, f.UploadFile)
}
