package user_controller

import (
	"github.com/gin-gonic/gin"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/global/variable"
	"hello/schoolMission/model"
	"hello/schoolMission/service/users/crud"
	"hello/schoolMission/service/users/token"
	"hello/schoolMission/utils/response"
	"time"
)

type UserController struct {

}

//  由于本项目骨架已经将表单验证器的字段(成员)绑定在上下文，因此可以按照 GetString()、GetInt64()、GetFloat64（）等快捷获取需要的数据类型，注意：相关键名规则：  前缀+验证器结构体中的 json 标签
// 当然也可以通过gin框架的上下文原始方法获取，例如： context.PostForm("user_name") 获取，这样获取的数据格式为文本，需要自己继续转换


// Register 用户注册
func (uc *UserController)Register(c *gin.Context)  {
	username := c.GetString(consts.ValidatorPrefix + "username")
	pass := c.GetString(consts.ValidatorPrefix + "pass")
	userIp := c.ClientIP()

	if crud.CreateCurdFactory().Register(username,pass,userIp){
		response.Success(c,consts.CurdStatusOkMsg, "")
	} else{
		response.Fail(c,consts.CurdRegisterFailCode,consts.CurdRegisterFailMsg,"")
	}
}



// Login 用户登录
func (uc *UserController) Login(c *gin.Context)  {
	userName := c.GetString(consts.ValidatorPrefix + "username")
	pass := c.GetString(consts.ValidatorPrefix + "pass")
	phone := c.GetString(consts.ValidatorPrefix + "phone")
	userModelFactory := model.CreateUserFactory("")
	userModel := userModelFactory.Login(userName,pass)

	if userModel != nil {
		userTokenFactory := token.CreateUserTokenFactory()
		if userToken, err := userTokenFactory.GenerateToken(userModel.Id,userModel.Username,
			variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt")); err == nil{
			if userTokenFactory.RecordLoginToken(userToken, c.ClientIP()){
				data := gin.H {
					"userId": userModel.Id,
					"user_name": userName,
					"phone": phone,
					"token": userToken,
					"updated_at": time.Now().Format(variable.DateFormat),
				}
				response.Success(c, consts.CurdStatusOkMsg, data)
				go userModel.UpdateUserLoginInfo(c.ClientIP(),userModel.Id)
				return
			}
		}
	}
	response.Fail(c,consts.CurdLoginFailCode,consts.CurdLoginFailMsg,"")
}

// RefreshToken 刷新用户token
func (uc *UserController) RefreshToken(c *gin.Context)  {
	oldToken := c.GetString(consts.ValidatorPrefix + "token")
	if newToken, ok := token.CreateUserTokenFactory().RefreshToken(oldToken, c.ClientIP()); ok{
		res := gin.H{
			"token": newToken,
		}
		response.Success(c, consts.CurdStatusOkMsg, res)
	} else {
		response.Fail(c,consts.CurdRefreshTokenFailCode,consts.CurdRefreshTokenFailMsg,"")
	}
}


// ShowByUserId 用户查询
func (uc *UserController) ShowByUserId(c *gin.Context)  {
	id := c.GetInt64(consts.ValidatorPrefix + "id")
	user, err := model.CreateUserFactory("").ShowOneItemById(id)
	if err != nil {
		response.Fail(c,consts.CurdSelectFailCode,consts.CurdSelectFailMsg, "")
	}
	response.Success(c,consts.CurdStatusOkMsg,user)
}


// ShowItemsByUsernamePattern 用户查询
func (uc *UserController) ShowItemsByUsernamePattern(c *gin.Context)  {
	usernamePattern := c.GetString(consts.ValidatorPrefix + "usernamePattern")
	page := c.GetInt64(consts.ValidatorPrefix + "page")
	limit := c.GetInt64(consts.ValidatorPrefix + "limit")
	user, err := model.CreateUserFactory("").ShowItemsByUsernamePattern(usernamePattern, page, limit)
	if err != nil {
		response.Fail(c,consts.CurdSelectFailCode,consts.CurdSelectFailMsg, "")
	}
	response.Success(c,consts.CurdStatusOkMsg,user)
}


// TODO: 不一定有必要实现
func (uc *UserController) ShowItemsByUserId(c *gin.Context) {

}





// Update 用户更新
func (uc *UserController) Update(c *gin.Context)  {
	userId := c.GetInt64(consts.ValidatorPrefix + "id")
	userName := c.GetString(consts.ValidatorPrefix + "username")
	pass := c.GetString(consts.ValidatorPrefix + "pass")
	phone := c.GetString(consts.ValidatorPrefix + "phone")
	identity := c.GetString(consts.ValidatorPrefix + "identity")
	birthdate := c.GetString(consts.ValidatorPrefix + "birthdate")
	userIp := c.ClientIP()
	//注意：这里没有实现权限控制逻辑，例如：超级管理管理员可以更新全部用户数据，普通用户只能修改自己的数据。目前只是验证了token有效、合法之后就可以进行后续操作
	// 实际使用请根据真是业务实现权限控制逻辑、再进行数据库操作
	if crud.CreateCurdFactory().Update(userId,userName,birthdate,pass,phone,identity,userIp) {
		response.Success(c, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(c, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}
}


