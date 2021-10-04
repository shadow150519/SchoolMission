package my_jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserId int64 `json:"id"`
	UserName string `json:"user_name"`
	Phone string `json:"phone"`
	jwt.StandardClaims
}
