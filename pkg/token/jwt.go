package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/towelong/lin-cms-go/pkg/response"
	"strings"
	"time"
)

type IToken interface {
	// GenerateToken 生成单令牌，需要指定令牌类型和过期时间
	GenerateToken(identity int, scope string, jwtType string, exp time.Time) string
	// GenerateAccessToken 生成access令牌
	GenerateAccessToken(identity int) string
	// GenerateRefreshToken 生成refresh令牌
	GenerateRefreshToken(identity int) string
	// GenerateTokens 生成双令牌
	GenerateTokens(identity int) Tokens
	// ParseToken 解析出令牌携带的信息
	ParseToken(tokenString string) (*Payload, error)
	// verifyToken 验证令牌是否有效以及是否过期
	verifyToken(tokenString string) (*Payload, error)
	// VerifyAccessToken 验证access令牌是否有效以及是否过期
	VerifyAccessToken(tokenString string) (*Payload, error)
	// VerifyRefreshToken 验证refresh令牌是否有效以及是否过期
	VerifyRefreshToken(tokenString string) (*Payload, error)
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DoubleJWT struct {
	Secret        string
	AccessExpire  time.Time
	RefreshExpire time.Time
}

// NewDoubleJWT
// accessExpire, refreshExpire 单位秒
func NewDoubleJWT(accessExpire, refreshExpire int64, secret string) IToken {
	access := getExpiredTime(accessExpire)
	refresh := getExpiredTime(refreshExpire)
	return &DoubleJWT{
		AccessExpire:  access,
		RefreshExpire: refresh,
		Secret:        secret,
	}
}

func (d *DoubleJWT) GenerateToken(identity int, scope string, jwtType string, exp time.Time) string {
	payload := NewPayload(identity, scope, jwtType, exp.Unix())
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := tokenClaims.SignedString([]byte(d.Secret))
	if err != nil {
		fmt.Printf("jwt generate error: %v", err)
	}
	return token
}

func (d *DoubleJWT) GenerateAccessToken(identity int) string {
	return d.GenerateToken(identity, Scope, AccessToken, d.AccessExpire)
}

func (d *DoubleJWT) GenerateRefreshToken(identity int) string {
	return d.GenerateToken(identity, Scope, RefreshToken, d.RefreshExpire)
}

func (d *DoubleJWT) GenerateTokens(identity int) Tokens {
	return Tokens{
		AccessToken:  d.GenerateAccessToken(identity),
		RefreshToken: d.GenerateRefreshToken(identity),
	}
}

// ParseToken 返回token携带的信息，并不做过期校验等
func (d DoubleJWT) ParseToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(d.Secret), nil
	})
	if token == nil {
		return nil, errors.New("jwt is nil")
	}
	return token.Claims.(*Payload), err
}

func (d DoubleJWT) verifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		payload := token.Claims.(*Payload)
		if payload.Scope != Scope {
			return nil, errors.New("jwt's scope is invalid")
		}
		return []byte(d.Secret), nil
	})
	// 处理token过期 以及 jwt 解析异常
	if err != nil {
		if err.(*jwt.ValidationError).Error() != "" {
			// 由于jwt库抛出的异常是比较笼统的
			// 所以想要把异常处理的细致些，需要通过异常信息来区分
			errMsg := err.Error()
			return nil, getErrorFromMessage(errMsg, token)
		}
		return nil, err
	}
	return token.Claims.(*Payload), nil
}

func (d DoubleJWT) VerifyAccessToken(tokenString string) (*Payload, error) {
	payload, err := d.verifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	if payload.Type == AccessToken {
		return payload, nil
	}
	return nil, response.NewResponse(10250)
}

func (d DoubleJWT) VerifyRefreshToken(tokenString string) (*Payload, error) {
	payload, err := d.verifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	if payload.Type == RefreshToken {
		return payload, nil
	}
	return nil, response.NewResponse(10250)
}

func getExpiredTime(n int64) time.Time {
	expire := time.Duration(n) * time.Second
	return time.Now().Add(expire)
}

// 拿到异常信息
func getErrorFromMessage(errMsg string, token *jwt.Token) error {
	var tokenType string
	resp := response.NewResponse(0)
	if token != nil {
		payload := token.Claims.(*Payload)
		tokenType = payload.Type
	}

	// 令牌过期
	if strings.Contains(errMsg, "expired") {
		if tokenType == AccessToken {
			resp.SetCode(10051)
			return resp
		}
		if tokenType == RefreshToken {
			resp.SetCode(10052)
			return resp
		}
		resp.SetCode(10050)
		return resp
	}
	// 令牌不合法
	if strings.Contains(errMsg, "scope") {
		resp.SetCode(10251)
		return resp
	}

	if strings.Contains(errMsg, "base64") || strings.Contains(errMsg, "signature") || strings.Contains(errMsg, "invalid number") {
		if tokenType == AccessToken {
			resp.SetCode(10041)
			return resp
		}
		if tokenType == RefreshToken {
			resp.SetCode(10042)
			return resp
		}
		resp.SetCode(10043)
		return resp
	}
	resp.SetCode(10013)
	return resp
}
