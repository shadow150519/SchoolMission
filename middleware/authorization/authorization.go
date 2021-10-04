package authorization

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hello/schoolMission/global/my_errors"
	"hello/schoolMission/global/variable"
	userstoken "hello/schoolMission/service/users/token"
	"hello/schoolMission/utils/response"
	"strings"
)

type HeaderParams struct {
	Authorization string `header:"Authorization"`
}

// CheckTokenAuth 检查Token权限
func CheckTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		headerParams := HeaderParams{}

		// 推荐使用 ShouldBindHeader 方式获取头参数
		if err := c.ShouldBindHeader(&headerParams); err != nil {
			variable.ZapLog.Error(my_errors.ErrorsValidatorBindParamsFail,zap.Error(err))
			c.Abort()
			return
		}

		if len(headerParams.Authorization) >= 20 {
			token := strings.Split(headerParams.Authorization, " ")
			if len(token) == 2 && len(token[1]) >= 20{
				tokenIsEffective := userstoken.CreateUserTokenFactory().IsEffective(token[1])
				if tokenIsEffective {
					if customeToken, err := userstoken.CreateUserTokenFactory().ParseToken(token[1]); err == nil{
						key := variable.ConfigYml.GetString("Token.BindContextKeyName")
						// token验证通过，同时绑定在请求上下文
						c.Set(key,customeToken)
					}
					c.Next()
				} else {
					response.ErrorTokenAuthFail(c)
					return
				}
			}
		} else{
			response.ErrorTokenAuthFail(c)
			return
		}
	}
}


