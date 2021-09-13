package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
	"net/http"
	"strconv"
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

func (admin *AdminAPI) ChangeUserPassword(ctx *gin.Context) {
	var password dto.ResetPasswordDTO
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "用户编号必须是正整数"))
		return
	}
	if err := ctx.ShouldBindJSON(&password); err != nil {
		ctx.Error(err)
		return
	}
	err := admin.UserService.ChangeUserPassword(id, password.NewPassword)
	if err != nil {
		ctx.Error(err)
		return
	}
	response.UpdatedVO(ctx)
}

func (admin *AdminAPI) DeleteUser(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "用户编号必须是正整数"))
		return
	}
	err := admin.UserService.DeleteUser(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	response.DeletedVO(ctx)
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
		admin.Auth.AdminRequired,
		admin.GetUsers,
	)
	adminRouter.LinPUT(
		"ChangeUserPassword",
		"/user/:id/password",
		adminRouter.Permission("修改用户密码", true),
		admin.Auth.AdminRequired,
		admin.ChangeUserPassword,
	)
	adminRouter.LinDELETE(
		"DeleteUser",
		"/user/:id",
		adminRouter.Permission("删除用户", true),
		admin.Auth.AdminRequired,
		admin.DeleteUser,
	)
}
