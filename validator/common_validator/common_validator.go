package common_validator

type BaseField struct {
	Username string `form:"username" json:"username"  binding:"required,min=1"` // 必填、对于文本,表示它的长度>=1
	Pass     string `form:"pass" json:"pass" binding:"required,min=6,max=20"`     //  密码为 必填，长度>=6
}

type Id struct {
	Id float64 `form:"id"  json:"id" binding:"required,min=1"`
}

type Page struct {
	Page  float64 `form:"page" json:"page" binding:"min=1"`   // 必填，页面值>=1
	Limit float64 `form:"limit" json:"limit" binding:"min=1"` // 必填，每页条数值>=1
}

