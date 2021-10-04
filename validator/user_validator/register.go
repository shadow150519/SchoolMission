package user_validator

import (
	"github.com/gin-gonic/gin"
	"hello/schoolMission/controller/user_controller"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/utils/response"
	"hello/schoolMission/validator/common_validator"
	"hello/schoolMission/validator/core/data_transfer"
)

type Register struct {
	common_validator.BaseField
	Phone string `from:"phone" json:"phone"`
	Identity string `from:"identity" json:"identity"`
	Birthdate string `form:"identity" json:"birthdate"`
}

// 特别注意: 表单参数验证器结构体的函数，绝对不能绑定在指针上
// 这部分代码项目启动会加载到容器中，如果绑定在指针上，一次请求之后，会造成容器中的代码段被污染

func (r Register) CheckParams(c *gin.Context){
	if err := c.ShouldBind(&r); err != nil{
		errs := gin.H{
			"tips": "UserRegister参数校验失败，参数不符合规定，username、pass、 长度有问题，不允许登录",
			"err":  err.Error(),
		}
		response.ErrorParam(c, errs)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(r, consts.ValidatorPrefix,c)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(c,"UserRegister表单验证器json化失败","")
	} else {
		// 验证完成调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&user_controller.UserController{}).Register(extraAddBindDataContext)
	}
}


