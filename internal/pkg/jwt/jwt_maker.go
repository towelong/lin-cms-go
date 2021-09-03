package jwt

import (
	"github.com/spf13/viper"
	jwt "github.com/towelong/lin-cms-go/pkg/token"
)

func NewJWTMaker() jwt.IToken {
	jwtToken := jwt.NewDoubleJWT(
		viper.GetInt64("lin.cms.tokenAccessExpire"),
		viper.GetInt64("lin.cms.tokenRefreshExpire"),
		viper.GetString("lin.cms.tokenSecret"),
	)
	return jwtToken
}
