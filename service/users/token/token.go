package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/global/my_errors"
	"hello/schoolMission/global/variable"
	"hello/schoolMission/model"
	my_jwt2 "hello/schoolMission/utils/my_jwt"
	"time"
)

type userToken struct {
	userJwt *my_jwt2.JwtSign
}

// CreateUserTokenFactory 创建userToken工厂
func CreateUserTokenFactory() *userToken{
	return &userToken{
		userJwt: my_jwt2.CreateMyJWT(variable.ConfigYml.GetString("Token.JwtTokenSignKey")),
	}
}

// GenerateToken 生成token
func (u *userToken) GenerateToken(userid int64, username string, expireAt int64)(tokens string, err error){
	// 根据实际业务自定义token需要包含的参数，生成token，注意: 用户密码请勿包含在token
	customClaims := my_jwt2.CustomClaims{
		UserId:             userid,
		UserName:       username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10, 		 // 生效开始时间
			ExpiresAt: time.Now().Unix() + expireAt, // 失效截止时间
		},
	}
	return u.userJwt.CreateToken(customClaims)
}

// RecordLoginToken 用户login成功, 记录用户token
func (u *userToken) RecordLoginToken(userToken, clientIp string) bool {
	if customClaims, err := u.userJwt.ParseToken(userToken);err == nil {
		userId := customClaims.UserId
		expiresAt := customClaims.ExpiresAt
		return model.CreateUserFactory("").OauthLoginToken(userId,userToken,expiresAt,clientIp)
	} else {
		return false
	}
}


// RefreshToken 刷新token的有效期 (默认+3600秒,参见常量配置项)
func (u *userToken) RefreshToken(oldToken, clientIp string)(newToken string, res bool){
	// 解析用户token的数据信息
	_, code := u.isNotExpired(oldToken)
	switch code {
		case consts.JwtTokenOK, consts.JwtTokenExpired:
			// 如果token已经过期，那么执行更新
			if newToken, err := u.userJwt.RefreshToken(oldToken,
				variable.ConfigYml.GetInt64("Token.JwtTokenRefreshExpireAt")); err == nil{
				if _ , err := u.userJwt.ParseToken(newToken); err == nil {
					return newToken, true
				}
			}
		case  consts.JwtTokenInvalid:
			variable.ZapLog.Error(my_errors.ErrorsTokenInvalid)
	}
	return "", false
}

// 判断token是否未过期
func (u *userToken) isNotExpired(token string)(*my_jwt2.CustomClaims, int)  {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {
		if time.Now().Unix() - customClaims.ExpiresAt < 0 {
			// token有效
			return customClaims, consts.JwtTokenOK
		} else {
			// 过期的token
			return customClaims, consts.JwtTokenExpired
		}
	} else {
		return nil, consts.JwtTokenInvalid
	}
}

// IsEffective 判断token是否有效 (未过期+数据库用户信息正常)
// TODO:不知道有什么用
func (u *userToken) IsEffective(token string) bool {
	_ , code := u.isNotExpired(token)
	return code == consts.JwtTokenOK
}

// 将token 解析为绑定时传递的参数
func (u *userToken) ParseToken(tokenStr string)(my_jwt2.CustomClaims,  error){
	if customClaims, err := u.userJwt.ParseToken(tokenStr); err == nil{
		return *customClaims, nil
	} else {
		return my_jwt2.CustomClaims{}, errors.New(my_errors.ErrorsParseTokenFail)
	}
}