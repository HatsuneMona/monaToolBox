package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"monaToolBox/app/userCenter/service"

	"monaToolBox/global"
	"monaToolBox/global/lock"
	"monaToolBox/global/response"
	"strconv"
	"time"
)

func JwtAuth(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authorization")
	if tokenStr == "" {
		response.ClaimsTokenFail(c)
		c.Abort()
		return
	}

	// Token 解析校验
	token, err := jwt.ParseWithClaims(
		tokenStr, &service.AdminUserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.Config.Jwt.Secret), nil
		},
	)
	// 如果验证不通过（token无法解码），或token在黑名单中
	// 则返回token错误
	if err != nil || service.JwtService.IsInBlacklist(tokenStr) {
		response.ClaimsTokenFail(c)
		c.Abort()
		return
	}
	// jwt 里的信息
	claims := token.Claims.(*service.AdminUserClaims)

	// 续签逻辑
	if claims.ExpiresAt.Unix()-time.Now().Unix() < global.Config.Jwt.RefreshGracePeriod {
		lock := lock.GetLock("refresh_token_lock", int64(global.Config.Jwt.JwtBlacklistGracePeriod))
		if lock.Lock() {
			userId, _ := strconv.Atoi(claims.UserId)
			err, user := service.UserService.GetUserInfoById(userId)
			if err != nil {
				global.Log.Error(err.Error())
				lock.Release()
			} else {
				tokenData, _ := service.JwtService.CreateToken(&user)
				c.Header("new-token", tokenData.AccessToken)
				c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
				_ = service.JwtService.JoinBlackList(token)
			}
		}
	}

	// Token 发布者校验
	c.Set("token", token)
	c.Set("userId", claims.UserId)
}
