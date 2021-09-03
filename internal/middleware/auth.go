package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
	"github.com/towelong/lin-cms-go/pkg/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "Bearer"
)

var AuthSet = wire.NewSet(wire.Struct(new(Auth), "*"))

type Auth struct {
	JWT          token.IToken
	UserService  service.IUserService
	GroupService service.IGroupService
}

func (a *Auth) LoginRequired(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		if err := a.mountUser(c); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.Next()
	} else {
		c.Next()
	}
}

func (a *Auth) GroupRequired(c *gin.Context) {
	if err := a.mountUser(c); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	user, _ := c.Get("currentUser")
	userId := user.(model.User).ID
	// admin直接通过
	admin, _ := a.UserService.IsAdmin(userId)
	if admin {
		c.Next()
	} else {
		meta, ok := c.Get("meta")
		if !ok {
			return
		}
		routeMeta := meta.(router.Meta)
		if !routeMeta.Mount {
			c.Next()
		} else {
			hasPermission := a.GroupService.GetUserHasPermission(userId, routeMeta)
			if !hasPermission {
				_ = c.Error(response.UnifyResponse(10001, c))
				c.Abort()
				return
			} else {
				c.Next()
			}
		}
	}
}

func (a *Auth) AdminRequired(c *gin.Context) {
	if err := a.mountUser(c); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	user, _ := c.Get("currentUser")
	currentUser := user.(model.User)
	admin, err := a.UserService.IsAdmin(currentUser.ID)
	if err != nil {
		_ = c.Error(response.UnifyResponse(10021, c))
		c.Abort()
		return
	}
	if admin {
		c.Next()
	} else {
		_ = c.Error(response.UnifyResponse(10001, c))
		c.Abort()
		return
	}
}

func (a *Auth) RefreshRequired(ctx *gin.Context) {
	refreshToken, tokenErr := getHeaderToken(ctx)
	if tokenErr != nil {
		ctx.Error(tokenErr)
		ctx.Abort()
		return
	}
	payload, err := a.JWT.VerifyRefreshToken(refreshToken)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
	} else {
		userId := payload.Identity
		user, errs := a.UserService.GetUserById(userId)
		if errs != nil {
			ctx.Error(response.UnifyResponse(10021, ctx))
			ctx.Abort()
		} else {
			ctx.Set("currentUser", user)
			ctx.Next()
		}
	}
}

func (a *Auth) mountUser(ctx *gin.Context) error {
	accessToken, tokenErr := getHeaderToken(ctx)
	if tokenErr != nil {
		return tokenErr
	}
	payload, err := a.JWT.VerifyAccessToken(accessToken)
	if err != nil {
		return err
	}
	// 校验用户
	userId := payload.Identity
	user, errs := a.UserService.GetUserById(userId)
	if errs != nil {
		return response.UnifyResponse(10021, ctx)
	}
	ctx.Set("currentUser", user)
	return nil
}

func getHeaderToken(ctx *gin.Context) (string, error) {
	authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
	if authorizationHeader == "" {
		return "", response.UnifyResponse(10012, ctx)
	}
	fields := strings.Fields(authorizationHeader)
	if fields[0] != AuthorizationTypeBearer {
		return "", response.UnifyResponse(10013, ctx)
	}
	tokenString := fields[1]
	return tokenString, nil
}
