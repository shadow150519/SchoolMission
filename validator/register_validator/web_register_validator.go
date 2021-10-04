package register_validator

import (
	"hello/schoolMission/core/container"
	"hello/schoolMission/global/consts"
	"hello/schoolMission/validator/upload_files"
	"hello/schoolMission/validator/user_validator"
)

// 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func WebRegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string
	// Users 模块表单验证器按照 key => value 形式注册在容器，方便路由模块中调用
	key = consts.ValidatorPrefix + "UsersRegister"
	containers.Set(key, user_validator.Register{})
	key = consts.ValidatorPrefix + "UsersLogin"
	containers.Set(key, user_validator.Login{})
	key = consts.ValidatorPrefix + "RefreshToken"
	containers.Set(key, user_validator.RefreshToken{})

	// Users基本操作（CURD）
	key = consts.ValidatorPrefix + "UsersShow"
	containers.Set(key, user_validator.ShowByUsername{})
	//key = consts.ValidatorPrefix + "UsersStore"
	//containers.Set(key, user_validator.Store{})
	key = consts.ValidatorPrefix + "UsersUpdate"
	containers.Set(key, user_validator.Update{})
	//key = consts.ValidatorPrefix + "UsersDestroy"
	//containers.Set(key, user_validator.{})

	// 文件上传
	key = consts.ValidatorPrefix + "UploadFiles"
	containers.Set(key, upload_files.UpFiles{})

	// Websocket 连接验证器
	//key = consts.ValidatorPrefix + "WebsocketConnect"
	//containers.Set(key, websocket.Connect{})
}
