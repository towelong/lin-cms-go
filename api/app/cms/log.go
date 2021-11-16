package cms

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/router"
)

type LogAPI struct {
	LogService service.ILogService
	Auth       middleware.Auth
}

func (l *LogAPI) GetUsers(ctx *gin.Context) {
	var dto dto.BasePage
	if err := ctx.ShouldBindQuery(&dto); err != nil {
		ctx.Error(err)
		return
	}
	users, _ := l.LogService.GetUsers(dto)
	ctx.JSON(http.StatusOK, users)
}

func (l *LogAPI) GetLogs(ctx *gin.Context) {
	var dto dto.BasePage
	if err := ctx.ShouldBindQuery(&dto); err != nil {
		ctx.Error(err)
		return
	}
	logs, _ := l.LogService.GetLogs(dto)
	ctx.JSON(http.StatusOK, logs)
}

func (l *LogAPI) SearchLogs(ctx *gin.Context) {
	var dto dto.SearchLogDTO
	if err := ctx.ShouldBindQuery(&dto); err != nil {
		ctx.Error(err)
		return
	}
	logs, _ := l.LogService.SearchLogs(dto)
	ctx.JSON(http.StatusOK, logs)
}

func (l *LogAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	logRouter := router.NewLinRouter("/log", "日志", routerGroup)
	logRouter.LinGET("GetLogs", "", logRouter.Permission("查询所有日志", true), l.Auth.GroupRequired, l.GetLogs)
	logRouter.LinGET("GetUsers", "/users", logRouter.Permission("查询日志记录的用户", true), l.Auth.GroupRequired, l.GetUsers)
	logRouter.LinGET("SearchLogs", "/search", logRouter.Permission("搜索日志", true), l.Auth.GroupRequired, l.SearchLogs)
}
