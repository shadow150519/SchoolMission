package my_jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"hello/schoolMission/global/my_errors"
	"time"
)

// JwtSign JWT验签 结构体
type JwtSign struct {
	SigningKey []byte
}

// CreateMyJWT 使用工厂创建一个 JWT 结构体
// 如果输入的signKe为空，则默认为goskeleton
func CreateMyJWT(signKey string)  *JwtSign {
	if len(signKey) <= 0 {
		signKey = "goskeleton"
	}
	return &JwtSign{SigningKey: []byte(signKey)}
}

// CreateToken 生成一个token
func (j *JwtSign) CreateToken(claims CustomClaims)(string, error)  {
	// 生成jwt的header,claims
	tokenPartA := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 添加密钥值，生成最后一部分
	return tokenPartA.SignedString(j.SigningKey)
}

// ParseToken 解析token
func (j *JwtSign) ParseToken(tokenString string)(*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if token == nil {
		return nil, errors.New(my_errors.ErrorsTokenInvalid)
	}
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError);ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0{
				return nil, errors.New(my_errors.ErrorsTokenMalFormed)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0{
				return nil, errors.New(my_errors.ErrorsTokenNotActiveYet)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// 如果TokenExpired，只是过期(格式都正确)，我们认为他是有效的，接下可以允许刷新操作
				token.Valid = true
				goto LABELHERE
			} else {
				return nil, errors.New(my_errors.ErrorsTokenInvalid)
			}
		}
	}
LABELHERE:
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
 	} else {
		 return nil, errors.New(my_errors.ErrorsTokenInvalid)
	}
}

// RefreshToken 更新token
func (j *JwtSign) RefreshToken(tokenString string, extraAddSeconds int64)(string, error) {
	if CustomClaims, err := j.ParseToken(tokenString); err == nil {
		CustomClaims.ExpiresAt = time.Now().Unix() + extraAddSeconds
		return j.CreateToken(*CustomClaims)
	} else {
		return "", err
	}
}


