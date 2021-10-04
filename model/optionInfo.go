package model

type OptionInfo struct {
	ID        uint `gorm:"primarykey"`
	SurveyID int `gorm:"not null"` // 关联调查问卷ID
	QuestionID int `gorm:"not null"` // 关联问题ID
	OptionName string `gorm:"not null"` // 选项描述
	OptionSort int `gorm:"not null"` // 选项排序，1-N
	OptionPicID int // 图片id
}
