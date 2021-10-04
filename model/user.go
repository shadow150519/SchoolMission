package model

import (
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"hello/schoolMission/global/variable"
	"hello/schoolMission/utils/md5_encrypt"
	my_jwt2 "hello/schoolMission/utils/my_jwt"
)

func CreateUserFactory(sqlType string) *User {
	return &User{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}


type User struct {
	BaseModel
	Birthdate string `gorm:"not null;default"`
	Password string  `gorm:"not null"`
	Phone string
	Token string
	Username string `gorm:"not null"`
	Identity string
	Level uint `gorm:"not null;default:2"` // 2 普通用户 1 管理员
	LastLoginIP string `gorm:"column:last_login_ip" json:"last_login_ip"`
}

// OauthLoginToken 记录用户登录生成的token，每次登录记录一次token
func (u *User) OauthLoginToken(userId int64, username string, expiresAt int64, clientIp string)bool  {
	token, err := my_jwt2.CreateMyJWT(variable.ConfigYml.GetString("Token.JwtTokenSignKey")).CreateToken(my_jwt2.CustomClaims{
		UserId:         userId,
		UserName:       username,
		Phone:          "",
		StandardClaims: jwt.StandardClaims{ExpiresAt:expiresAt},
	})
	if err != nil {
		variable.ZapLog.Error("生成token错误:", zap.Error(err))
		return false
	}
	sql := "UPDATE `user_tokens` SET token=? where id=?"
	result := u.Raw(sql,token,userId)
	if result.Error != nil {
		variable.ZapLog.Error("设置登录token失败",zap.Error(result.Error))
		return false
	}
	return true

}

// RefreshToken 用户刷新token
func (u *User) RefreshToken(userId int64, newToken, clientIp string)bool  {
	sql := "UPDATE users SET token=? WHERE userId=?"
	if u.Raw(sql,newToken,userId).Error == nil{
		go u.UpdateUserLoginInfo(clientIp, userId)
		return true
	}
	return false
}

// IsUserExist 根据UserName判断用户是否存在
func (u *User) IsUserExist(userName string) bool {
	sql := "select 1 from users where username = ?"
	result := u.Raw(sql,userName)
	return result.RowsAffected != 0
}

// UpdateUserLoginInfo 更新用户最近一次登录ip
func (u *User) UpdateUserLoginInfo(lastLoginIp string, userId int64)  {
	sql := "UPDATE users SET last_login_ip=?, where id=?"
	u.Exec(sql, lastLoginIp, userId)
}


// Register 用户注册
func (u *User) Register(userName, password, userIp string)bool  {
	sql := "INSERT  INTO users(username,password,last_login_ip) SELECT ?,?,? FROM DUAL " +
		" WHERE NOT EXISTS (SELECT 1  FROM users WHERE  username=?)"

	result := u.Raw(sql, userName, password, userIp, userName)

	return result.RowsAffected > 0

}

// Login 用户登录
func (u *User)Login(username, password string)*User {
	sql := "select `id`, `username`, `pass`, `phone`, `birthdate`, `identity` from users where username = ?"
	result := u.Raw(sql, username).First(u)
	if result.Error == nil {
		// 帐号密码验证成功
		if len(u.Password) > 0 && (u.Password == md5_encrypt.Base64MD5(password)){
			return u
		}
	} else {
		variable.ZapLog.Error("根据账号查询单条记录出错:", zap.Error(result.Error))
	}
	return nil
}

// ResetToken 当用户更改密码后，所有的token都失效，必须重新登陆
// TODO: 回头搞
func (u *User) ResetToken(id int64, newPassword, clientIP string)bool {
	// 如果用户新旧密码一致则不需要重新登录,返回true
	userItem, err := u.ShowOneItemById(id)
	if userItem != nil && err == nil && userItem.Password == newPassword {
		return true
	} else if userItem != nil {
		panic("ResetToken")
	}
	return false
}

// ShowOneItemById 根据用户名查询一条信息
func (u *User) ShowOneItemById(userId int64)(*User, error){
	sql := "SELECT `id`, `username`, `phone`, `birthdate`, `identity` from users where id = ? "
	result := u.Raw(sql, userId).First(u)
	if result.Error == nil {
		return u, nil
 	} else {
		 return nil, result.Error
	}
}

// ShowItemsByUsernamePattern ShowItemsByUsername 根据用户名模糊查询信息
func (u *User) ShowItemsByUsernamePattern(usernamePattern string, page, limit int64)([]*User, error)  {
	users := []*User{}
	sql := "SELECT `id`, `username`, `phone`, `birthdate`, `identity` from users " +
		"where username like ? LIMIT ?, ?"
	result := u.Raw(sql, "%"+usernamePattern+"%", (page - 1) * 10, limit).Find(users)
	if result.Error == nil {
		return users, nil
	} else {
		return nil, result.Error
	}
}


// Update 更新用户信息
func (u *User)Update(id int64, username string, password string, identity string, phone string, clientIp string)bool{
	sql := "UPDATE users (username, password, identity, phone, clientIp) VALUES(?,?,?,?,?) WHERE id = ?"
	if u.Exec(sql,username,password,identity,phone,clientIp,id).RowsAffected >= 0 {
		if u.ResetToken(id, password, clientIp) {
			return true
		}
	}
	return false
}

