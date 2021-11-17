package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
	"github.com/towelong/lin-cms-go/pkg/token"
	"net/http"
)

type UserAPI struct {
	JWT          token.IToken
	UserService  service.IUserService
	GroupService service.IGroupService
	Auth         middleware.Auth
}

func (u *UserAPI) Register(ctx *gin.Context) {
	var register dto.RegisterDTO
	if err := ctx.ShouldBindJSON(&register); err != nil {
		ctx.Error(err)
		return
	}
	if err := u.UserService.CreateUser(register); err != nil {
		ctx.Error(err)
		return
	}
	response.CreatedVO(ctx, 11)
}

func (u *UserAPI) Login(ctx *gin.Context) {
	var login dto.UserLoginDTO
	validatorErr := ctx.ShouldBindJSON(&login)
	if validatorErr != nil {
		ctx.Error(validatorErr)
		return
	}
	user, userErr := u.UserService.VerifyUser(login.Username, login.Password)
	if userErr != nil {
		ctx.Error(userErr)
		return
	}
	tokens := u.JWT.GenerateTokens(user.ID)
	ctx.JSON(http.StatusOK, tokens)
}

func (u *UserAPI) UpdatePassword(ctx *gin.Context) {
	var passwordDTO dto.ChangePasswordDTO
	if err := ctx.ShouldBindJSON(&passwordDTO); err != nil {
		ctx.Error(err)
		return
	}
	user, _ := ctx.Get("currentUser")
	currentUser := user.(model.User)
	if err := u.UserService.ChangePassword(currentUser.ID, passwordDTO); err != nil {
		ctx.Error(err)
		return
	}
	response.UpdatedVO(ctx, 4)
}

func (u *UserAPI) Update(ctx *gin.Context) {
	var userInfoDTO dto.UpdateInfoDTO
	if err := ctx.ShouldBindJSON(&userInfoDTO); err != nil {
		ctx.Error(err)
		return
	}
	user, _ := ctx.Get("currentUser")
	currentUser := user.(model.User)
	if err := u.UserService.UpdateProfile(currentUser.ID, userInfoDTO); err != nil {
		ctx.Error(err)
		return
	}
	response.UpdatedVO(ctx, 6)
}

func (u *UserAPI) RefreshToken(ctx *gin.Context) {
	user, _ := ctx.Get("currentUser")
	currentUser := user.(model.User)
	tokens := u.JWT.GenerateTokens(currentUser.ID)
	ctx.JSON(http.StatusOK, tokens)
}

func (u *UserAPI) GetInformation(ctx *gin.Context) {
	user, _ := ctx.Get("currentUser")
	currentUser := user.(model.User)
	groups := u.UserService.GetUserGroupByUserId(currentUser.ID)
	var userInfo vo.UserInfo
	copier.Copy(&userInfo, &currentUser)
	copier.CopyWithOption(&userInfo.Groups, &groups, copier.Option{IgnoreEmpty: true})
	if userInfo.Groups == nil {
		userInfo.Groups = make([]vo.Group, 0)
	}
	ctx.JSON(http.StatusOK, userInfo)
}

func (u *UserAPI) GetPermissions(ctx *gin.Context) {
	user, _ := ctx.Get("currentUser")
	currentUser := user.(model.User)
	info, err := u.UserService.GetUserPermissionsInfo(currentUser.ID)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, info)
}

func (u UserAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	userRouter := router.NewLinRouter("/user", "用户", routerGroup)
	userRouter.POST("/register", u.Auth.AdminRequired, u.Register)
	userRouter.PUT("/change_password", u.Auth.LoginRequired, u.UpdatePassword)
	userRouter.POST("/login", u.Login)
	userRouter.PUT("", u.Auth.LoginRequired, u.Update)
	userRouter.GET("/refresh", u.Auth.RefreshRequired, u.RefreshToken)
	userRouter.GET("/information", u.Auth.LoginRequired, u.GetInformation)
	userRouter.GET("/permissions", u.Auth.LoginRequired, u.GetPermissions)
}
