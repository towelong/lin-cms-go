package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/middleware"

	"net/http"
	"strconv"

	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
)

type AdminAPI struct {
	PermissionService service.IPermissionService
	UserService       service.IUserService
	GroupService      service.IGroupService
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

func (admin *AdminAPI) GetGroups(ctx *gin.Context) {
	var page dto.BasePage
	if err := ctx.ShouldBindQuery(&page); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, admin.GroupService.GetPageGroups(page))
}

func (admin *AdminAPI) GetAllGroups(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, admin.GroupService.GetAllGroups())
}

func (admin *AdminAPI) GetGroup(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "id必须是正整数"))
		return
	}
	groupInfo, err := admin.GroupService.GetGroupById(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, groupInfo)
}

func (admin *AdminAPI) CreateGroup(ctx *gin.Context) {
	var groupDTO dto.NewGroupDTO
	if err := ctx.ShouldBindJSON(&groupDTO); err != nil {
		ctx.Error(err)
		return
	}
	err := admin.GroupService.CreateGroup(groupDTO)
	if err != nil {
		ctx.Error(err)
		return
	}
	response.CreatedVO(ctx, 15)
}

func (admin *AdminAPI) UpdateGroup(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "id必须是正整数"))
		return
	}
	var groupDTO dto.UpdateGroupDTO
	if err := ctx.ShouldBindJSON(&groupDTO); err != nil {
		ctx.Error(err)
		return
	}
	if err := admin.GroupService.UpdateGroup(id, groupDTO); err != nil {
		ctx.Error(err)
		return
	}
	response.UpdatedVO(ctx, 7)
}

func (admin *AdminAPI) DeleteGroup(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "id必须是正整数"))
		return
	}
	if err := admin.GroupService.DeleteGroup(id); err != nil {
		ctx.Error(err)
		return
	}
	response.DeletedVO(ctx, 8)
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
	response.UpdatedVO(ctx, 4)
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
	response.DeletedVO(ctx, 5)
}

func (admin *AdminAPI) UpdateUser(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "用户编号必须是正整数"))
		return
	}
	var userInfoDTO dto.UpdateGroupsDTO
	if err := ctx.ShouldBindJSON(&userInfoDTO); err != nil {
		ctx.Error(err)
		return
	}
	if err := admin.UserService.UpdateUserInfo(id, userInfoDTO); err != nil {
		ctx.Error(err)
		return
	}
	response.UpdatedVO(ctx, 6)
}

func (admin *AdminAPI) DispatchPermission(ctx *gin.Context) {
	var dispatchDTO dto.DispatchPermissionDTO
	if err := ctx.ShouldBindJSON(&dispatchDTO); err != nil {
		ctx.Error(err)
		return
	}
	if err := admin.PermissionService.DispatchPermission(dispatchDTO); err != nil {
		ctx.Error(err)
		return
	}
	response.CreatedVO(ctx, 9)
}

func (admin *AdminAPI) DispatchPermissions(ctx *gin.Context) {
	var dispatchDTO dto.DispatchPermissionsDTO
	if err := ctx.ShouldBindJSON(&dispatchDTO); err != nil {
		ctx.Error(err)
		return
	}
	if err := admin.PermissionService.DispatchPermissions(dispatchDTO); err != nil {
		ctx.Error(err)
		return
	}
	response.CreatedVO(ctx, 9)
}

func (admin *AdminAPI) RemovePermissions(ctx *gin.Context) {
	var dispatchDTO dto.DispatchPermissionsDTO
	if err := ctx.ShouldBindJSON(&dispatchDTO); err != nil {
		ctx.Error(err)
		return
	}
	if err := admin.PermissionService.RemovePermissions(dispatchDTO); err != nil {
		ctx.Error(err)
		return
	}
	response.DeletedVO(ctx, 10)
}

func (admin *AdminAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	adminRouter := router.NewLinRouter("/admin", "管理员", routerGroup)
	adminRouter.LinGET(
		"GetAllPermissions",
		"/permission",
		adminRouter.Permission("查询所有可分配的权限", false),
		admin.Auth.AdminRequired,
		admin.GetAllPermissions,
	)
	adminRouter.LinGET(
		"GetUsers",
		"/users",
		adminRouter.Permission("查询所有用户", false),
		admin.Auth.AdminRequired,
		admin.GetUsers,
	)
	adminRouter.LinPUT(
		"ChangeUserPassword",
		"/user/:id/password",
		adminRouter.Permission("修改用户密码", false),
		admin.Auth.AdminRequired,
		admin.ChangeUserPassword,
	)
	adminRouter.LinDELETE(
		"DeleteUser",
		"/user/:id",
		adminRouter.Permission("删除用户", false),
		admin.Auth.AdminRequired,
		admin.DeleteUser,
	)
	adminRouter.LinPUT(
		"UpdateUser",
		"/user/:id",
		adminRouter.Permission("管理员更新用户信息", false),
		admin.Auth.AdminRequired,
		admin.UpdateUser,
	)
	adminRouter.LinGET(
		"GetGroups",
		"/group",
		adminRouter.Permission("查询所有权限组及其权限", false),
		admin.Auth.AdminRequired,
		admin.GetGroups,
	)
	adminRouter.LinGET(
		"GetAllGroups",
		"/group/all",
		adminRouter.Permission("查询所有权限组", false),
		admin.Auth.AdminRequired,
		admin.GetAllGroups,
	)
	adminRouter.LinPOST(
		"CreateGroup",
		"/group",
		adminRouter.Permission("新建权限组", false),
		admin.Auth.AdminRequired,
		admin.CreateGroup,
	)
	adminRouter.LinGET(
		"GetGroup",
		"/group/:id",
		adminRouter.Permission("查询一个权限组及其权限", false),
		admin.Auth.AdminRequired,
		admin.GetGroup,
	)
	adminRouter.LinPUT(
		"UpdateGroup",
		"/group/:id",
		adminRouter.Permission("更新一个权限组", false),
		admin.Auth.AdminRequired,
		admin.UpdateGroup,
	)
	adminRouter.LinDELETE(
		"DeleteGroup",
		"/group/:id",
		adminRouter.Permission("删除一个权限组", false),
		admin.Auth.AdminRequired,
		admin.DeleteGroup,
	)
	adminRouter.LinPOST(
		"DispatchPermission",
		"/permission/dispatch",
		adminRouter.Permission("分配单个权限", false),
		admin.Auth.AdminRequired,
		admin.DispatchPermission,
	)
	adminRouter.LinPOST(
		"DispatchPermissions",
		"/permission/dispatch/batch",
		adminRouter.Permission("分配多个权限", false),
		admin.Auth.AdminRequired,
		admin.DispatchPermissions,
	)
	adminRouter.LinDELETE(
		"RemovePermissions",
		"/permission/remove",
		adminRouter.Permission("删除多个权限", false),
		admin.Auth.AdminRequired,
		admin.RemovePermissions,
	)
}
