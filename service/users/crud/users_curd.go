package crud

import (
	"hello/schoolMission/model"
	"hello/schoolMission/utils/md5_encrypt"
)

func CreateCurdFactory() *UsersCurd{
	return &UsersCurd{}
}
type UsersCurd struct {
}


func (u *UsersCurd) Register(userName, pass, userIp string)bool {
	pass = md5_encrypt.Base64MD5(pass)
	// sqlType 为空时使用配置文件里的sqlType
	return model.CreateUserFactory("").Register(userName, pass, userIp)
}

func (u *UsersCurd) Update(id int64, username, birthdate, password, phone, identity, lastLoginIp string)bool{
	password = md5_encrypt.Base64MD5(password)
	return model.CreateUserFactory("").Update(id,username,password,identity,phone,lastLoginIp)
}


