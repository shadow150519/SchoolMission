package user_validator

import (
	"github.com/gin-gonic/gin"
	"hello/schoolMission/controller/user_controller"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/utils/response"
	"strings"
)

type RefreshToken struct {
	Authorization string `json:"token" header:"Authorization" binding:"required,min=20"`
}

func (rft RefreshToken) CheckParams(c *gin.Context)  {
	if err := c.ShouldBindHeader(&rft); err != nil {
		errs := gin.H{
			"tips": "Token参数校验失败，参数不符合规定，token 长度>=20",
			"err":  err.Error(),
		}
		response.ErrorParam(c, errs)
		return
	}
	token := strings.Split(rft.Authorization, " ")
	// bear <token>
	if len(token) == 2 {
		c.Set(consts.ValidatorPrefix+"token", token[1])
		(&user_controller.UserController{}).RefreshToken(c)
	} else {
		errs := gin.H{
			"tips": "Token不合法，token请放置在header头部分，按照按=>键提交，例如：Authorization：Bearer 你的实际token....",
		}
		response.Fail(c, consts.JwtTokenFormatErrCode, consts.JwtTokenFormatErrMsg, errs)
	}
}
