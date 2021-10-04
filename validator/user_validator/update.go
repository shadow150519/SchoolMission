package user_validator

import (
	"github.com/gin-gonic/gin"
	"hello/schoolMission/controller/user_controller"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/utils/response"
	"hello/schoolMission/validator/common_validator"
	"hello/schoolMission/validator/core/data_transfer"
)

type Update struct {
	common_validator.BaseField
	common_validator.Id
	Phone string `form:"phone" json:"phone" binding:"required,len=11"`
	Identity string `form:"indentity" json:"identity" binding:"required"`
}

func (u Update) CheckParams(c *gin.Context)  {
	if err := c.ShouldBind(&u); err != nil {
		errs := gin.H{
			"tips": "UserUpdate，参数校验失败，请检查id(>0),user_name(>=1)、pass(>=6)、real_name(>=2)、phone长度(=11)",
			"err":  err.Error(),
		}
		response.ErrorParam(c, errs)
		return
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文，便于下一步（控制器）可以直接通过 context.Get(键) 获取相关值
	extraAddBindDataContext := data_transfer.DataAddContext(u, consts.ValidatorPrefix, c)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(c, "UserUpdate表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&user_controller.UserController{}).Update(extraAddBindDataContext)
	}
}
