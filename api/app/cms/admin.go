package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
	"net/http"
)

type AdminAPI struct {
	PermissionService service.IPermissionService
	UserService       service.IUserService
	Auth              middleware.Auth
}

func (admin *AdminAPI) GetAllPermissions(ctx *gin.Context) {
	permissions, err := admin.PermissionService.GetStructPermissions()
	if err != nil {
		response.NotFound(ctx)
	} else {
		ctx.JSON(200, permissions)
	}
}

func (admin *AdminAPI) GetUsers(ctx *gin.Context) {
	var queryUser dto.QueryUserDTO
	if err := ctx.ShouldBindQuery(&queryUser); err != nil {
		ctx.Error(err)
		return
	}
	userPage, err := admin.UserService.GetUserPageByGroupId(queryUser.GroupId, queryUser.Page, queryUser.Count)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, userPage)
}

func (admin *AdminAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	adminRouter := router.NewLinRouter("/admin", "管理员", routerGroup)
	adminRouter.LinGET(
		"GetAllPermissions",
		"/permission",
		adminRouter.Permission("查询所有可分配的权限", true),
		admin.Auth.AdminRequired,
		admin.GetAllPermissions,
	)
	adminRouter.LinGET(
		"GetUsers",
		"/users",
		adminRouter.Permission("查询所有用户", true),
		admin.Auth.GroupRequired,
		admin.GetUsers,
	)
}
