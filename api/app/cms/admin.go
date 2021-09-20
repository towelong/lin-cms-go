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
	response.CreatedVO(ctx)
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
	response.UpdatedVO(ctx)
}

func (admin *AdminAPI) DeleteGroup(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "id必须是正整数"))
		return
	}
	if err := admin.GroupService.DeleteGroup(id);err!=nil {
		ctx.Error(err)
		return
	}
	response.DeletedVO(ctx)
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

	adminRouter.LinGET(
		"GetGroups",
		"/group",
		adminRouter.Permission("查询所有权限组及其权限", true),
		admin.Auth.AdminRequired,
		admin.GetGroups,
	)
	adminRouter.LinGET(
		"GetAllGroups",
		"/group/all",
		adminRouter.Permission("查询所有权限组", true),
		admin.Auth.AdminRequired,
		admin.GetAllGroups,
	)
	adminRouter.LinPOST(
		"CreateGroup",
		"/group",
		adminRouter.Permission("新建权限组", true),
		admin.Auth.AdminRequired,
		admin.CreateGroup,
	)
	adminRouter.LinGET(
		"GetGroup",
		"/group/:id",
		adminRouter.Permission("查询一个权限组及其权限", true),
		admin.Auth.AdminRequired,
		admin.GetGroup,
	)
	adminRouter.LinPUT(
		"UpdateGroup",
		"/group/:id",
		adminRouter.Permission("更新一个权限组", true),
		admin.Auth.AdminRequired,
		admin.UpdateGroup,
	)
	adminRouter.LinDELETE(
		"DeleteGroup",
		"/group/:id",
		adminRouter.Permission("删除一个权限组", true),
		admin.Auth.AdminRequired,
		admin.DeleteGroup,
	)
}
