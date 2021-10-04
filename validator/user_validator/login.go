package user_validator

import (
	"github.com/gin-gonic/gin"
	"hello/schoolMission/controller/user_controller"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/utils/response"
	"hello/schoolMission/validator/common_validator"
	"hello/schoolMission/validator/core/data_transfer"
)

type Login struct {
	common_validator.BaseField
}

func (l Login) CheckParams(c *gin.Context)  {
	if err := c.ShouldBind(&l); err != nil {
		errs := gin.H{
			"tips": "UserRegister参数校验失败，参数不符合规定，user_name、pass、 长度有问题，不允许登录",
			"err":  err.Error(),
		}
		response.ErrorParam(c, errs)
		return
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文，便于下一步（控制器）可以直接通过 context.Get(键) 获取相关值
	extraAddBindDataContext := data_transfer.DataAddContext(l, consts.ValidatorPrefix, c)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(c, "userLogin表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&user_controller.UserController{}).Login(extraAddBindDataContext)
	}
}
