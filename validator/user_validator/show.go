package user_validator

import (
	"github.com/gin-gonic/gin"
	"hello/schoolMission/controller/user_controller"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/utils/response"
	"hello/schoolMission/validator/common_validator"
	"hello/schoolMission/validator/core/data_transfer"
)

type ShowByUsername struct {
	Username string `form:"username" json:"username" binding:"required,min=1"`
	common_validator.Page
}

func (sbu ShowByUsername) CheckParams(c *gin.Context)  {
	// 1.基本的验证规则没有通过
	if err := c.ShouldBind(&sbu); err != nil{
		errs := gin.H{
			"tips": "UserShow参数校验失败，参数不符合规定，user_name（长度>0）、page的值(>0)、limits的值（>0)",
			"err":  err.Error(),
		}
		response.ErrorParam(c,errs)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(sbu, consts.ValidatorPrefix,c)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(c, "UserShow表单验证器json化失败", "")
	} else {
		(&user_controller.UserController{}).ShowItemsByUsernamePattern(extraAddBindDataContext)
	}
}

