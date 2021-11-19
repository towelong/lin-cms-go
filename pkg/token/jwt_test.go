package token

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/towelong/lin-cms-go/internal/config"
	"github.com/towelong/lin-cms-go/pkg/response"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	t.Run("测试token生成", func(t *testing.T) {
		var accessExpired int64 = 60
		jwt := NewDoubleJWT(accessExpired, 3600, "xxxx")
		token := jwt.GenerateAccessToken(1)
		refreshToken := jwt.GenerateRefreshToken(1)
		tokens := jwt.GenerateTokens(1)

		now := time.Now()
		expired := now.Add(time.Duration(accessExpired) * time.Second).Unix()

		require.NotEmpty(t, token)
		require.NotEmpty(t, tokens)
		payload, err := jwt.verifyToken(token)
		accessPayload, err2 := jwt.VerifyAccessToken(token)
		refreshPayload, err3 := jwt.VerifyRefreshToken(refreshToken)
		require.NoError(t, err)
		require.NoError(t, err2)
		require.NoError(t, err3)

		require.NotZero(t, payload.Identity)
		require.Equal(t, payload.Exp, expired)
		require.Equal(t, AccessToken, accessPayload.Type)
		require.Equal(t, RefreshToken, refreshPayload.Type)

		require.Nil(t, err2)
		require.Nil(t, err3)
	})
}

func TestInvalidToken(t *testing.T) {
	// case1: 测试过期token
	var accessExpired int64 = -60
	var refreshExpired int64 = -60
	jwt := NewDoubleJWT(accessExpired, refreshExpired, "xxxx")
	t.Run("case1: 测试过期token", func(t *testing.T) {
		token := jwt.GenerateAccessToken(1)
		refreshToken := jwt.GenerateRefreshToken(1)
		require.NotEmpty(t, token)
		require.NotEmpty(t, refreshToken)
		payload, err := jwt.verifyToken(token)
		payloadRefresh, refreshErr := jwt.verifyToken(refreshToken)
		require.Error(t, err)
		require.Error(t, refreshErr)
		require.EqualError(t, err, response.NewResponse(10051).Error())
		require.EqualError(t, refreshErr, response.NewResponse(10052).Error())
		require.Nil(t, payload)
		require.Nil(t, payloadRefresh)
	})

	// case2: 测试不合法token
	t.Run("case2:测试不合法token", func(t *testing.T) {
		token := jwt.GenerateAccessToken(2)
		rToken := jwt.GenerateRefreshToken(2)
		aPayload, err4 := jwt.verifyToken(token + "aa")
		rPayload, err5 := jwt.verifyToken(rToken + "aa")
		require.Error(t, err4)
		require.Error(t, err5)
		require.EqualError(t, err4, response.NewResponse(10041).Error())
		require.EqualError(t, err5, response.NewResponse(10042).Error())
		require.Nil(t, aPayload)
		require.Nil(t, rPayload)

		payload2, err2 := jwt.verifyToken("11111")
		payload, err := jwt.VerifyAccessToken("666")
		payload3, err3 := jwt.VerifyRefreshToken("666")
		require.Error(t, err2)
		require.Error(t, err)
		require.Error(t, err3)
		require.EqualError(t, err2, response.NewResponse(10043).Error())
		require.EqualError(t, err, response.NewResponse(10043).Error())
		require.EqualError(t, err3, response.NewResponse(10043).Error())
		require.Nil(t, payload2)
		require.Nil(t, payload)
		require.Nil(t, payload3)
	})

	// case2-1: 测试scope不合法的token
	t.Run("case2-1:测试scope不合法的token", func(t *testing.T) {
		generateToken := jwt.GenerateToken(1, "welong", AccessToken, getExpiredTime(60))
		require.NotEmpty(t, generateToken)
		verifyToken, err3 := jwt.verifyToken(generateToken)
		require.Error(t, err3)
		require.EqualError(t, err3, response.NewResponse(10251).Error())
		require.Nil(t, verifyToken)
	})

	// case2-2: 测试既不是accessToken 也不是 refreshToken
	t.Run("case2-2:测试既不是accessToken,也不是 refreshToken", func(t *testing.T) {
		token := jwt.GenerateToken(1, Scope, "no", getExpiredTime(30))
		token2 := jwt.GenerateToken(1, Scope, "no", getExpiredTime(-30))
		require.NotEmpty(t, token)
		require.NotEmpty(t, token2)

		access, err := jwt.VerifyAccessToken(token)
		refresh, err2 := jwt.VerifyRefreshToken(token)
		expiredToken, err3 := jwt.verifyToken(token2)
		require.Error(t, err)
		require.Error(t, err2)
		require.Error(t, err3)
		require.EqualError(t, err, response.NewResponse(10250).Error())
		require.EqualError(t, err2, response.NewResponse(10250).Error())
		require.EqualError(t, err3, response.NewResponse(10050).Error())
		require.Nil(t, access)
		require.Nil(t, refresh)
		require.Nil(t, expiredToken)
	})

	t.Run("case2-3: 测试token携带的信息", func(t *testing.T) {
		doubleJWT := NewDoubleJWT(60, 60, "xxx")
		token := doubleJWT.GenerateAccessToken(1)
		payload, err := doubleJWT.ParseToken(token)
		require.NoError(t, err)
		require.NotEmpty(t, payload)
		require.Equal(t, uint(1), payload.Identity)
	})
}

func TestDoubleJWT_GenerateToken(t *testing.T) {
	conf := config.LoadConfig()
	jwt := NewDoubleJWT(conf.Lin.CMS.TokenAccessExpire, conf.Lin.CMS.TokenRefreshExpire, conf.Lin.CMS.TokenSecret)
	token := jwt.GenerateAccessToken(1)
	//assert.Equal(t, "ss", token)
	assert.NotEmpty(t, token)
}
