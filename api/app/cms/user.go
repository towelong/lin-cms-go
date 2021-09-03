package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/router"
	"github.com/towelong/lin-cms-go/pkg/token"
	"net/http"
)

type UserAPI struct {
	JWT         token.IToken
	UserService service.IUserService
	Auth        middleware.Auth
}

func (u *UserAPI) Register(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"hello": ctx,
	})
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

func (u UserAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	userRouter := router.NewLinRouter("/user", "用户", routerGroup)
	userRouter.LinPOST(
		"Register",
		"/register",
		userRouter.Permission("用户注册", true),
		u.Auth.AdminRequired,
		u.Register,
	)
	userRouter.POST("/login", u.Login)
}
