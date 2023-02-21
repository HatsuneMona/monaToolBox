package service

import (
	"context"
	jwt "github.com/golang-jwt/jwt/v4"
	"monaToolBox/global"
	"monaToolBox/utils"
	"strconv"
	"time"
)

type jwtService struct {
}

var JwtService = new(jwtService)

type JwtUser interface {
	GetUid() string
}

type AdminUserClaims struct {
	jwt.RegisteredClaims
	// JWT 官方规定的七个字段
	// iss (issuer)：签发人
	// exp (expiration time)：过期时间
	// sub (subject)：主题
	// aud (audience)：受众
	// nbf (Not Before)：生效时间
	// iat (Issued At)：签发时间
	// jti (JWT ID)：编号
	UserId string `json:"user_id"`
}

type claimsToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// CreateToken 为用户生成 jwt token
func (jwtService *jwtService) CreateToken(user JwtUser) (token claimsToken, err error) {
	tokenData := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		AdminUserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(global.Config.Jwt.JwtTtl))), // 过期时间
				NotBefore: jwt.NewNumericDate(time.Now()),                                                            // 生效时间
				IssuedAt:  jwt.NewNumericDate(time.Now()),                                                            // 签发时间
			},
			UserId: user.GetUid(),
		},
	)

	tokenStr, err := tokenData.SignedString([]byte(global.Config.Jwt.Secret))

	token = claimsToken{
		AccessToken: tokenStr,
		ExpiresIn:   global.Config.Jwt.JwtTtl,
	}

	return
}

// getBlackListKey 生成jwt的唯一key
func (jwtService *jwtService) getBlackListKey(tokenStr string) string {
	return "jwt_black_list:" + utils.MD5([]byte(tokenStr))
}

// JoinBlackList token 加入黑名单
func (jwtService *jwtService) JoinBlackList(token *jwt.Token) (err error) {
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*AdminUserClaims).ExpiresAt.Unix()-nowUnix) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	err = global.Redis.SetNX(context.Background(), jwtService.getBlackListKey(token.Raw), nowUnix, timer).Err()
	return
}

// IsInBlacklist token 是否在黑名单中
func (jwtService *jwtService) IsInBlacklist(tokenStr string) bool {
	joinUnixStr, err := global.Redis.Get(context.Background(), jwtService.getBlackListKey(tokenStr)).Result()
	joinBlackUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if joinUnixStr == "" || err != nil {
		return false
	}
	// JwtBlacklistGracePeriod 为黑名单宽限时间，避免并发请求失效
	if time.Now().Unix()-joinBlackUnix < int64(global.Config.Jwt.JwtBlacklistGracePeriod) {
		return false
	}
	return true
}
