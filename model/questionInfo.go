package model

type QuestionInfo struct {
	ID        uint `gorm:"primarykey"`
	SurveyID int `gorm:"not null"` // 关联调查问卷ID
	QuestionType byte `gorm:"not null"` // 题目类型，1 单选 2多选 3填空 4文件
	QuestionName string `gorm:"not null"` // 问题主题
	QuestionDescription string `gorm:"not null"` // 问题描述
	QuestionSort int `gorm:"not null"` // 问题在问卷中是第几题，从1开始
	RequiredFlag byte `gorm:"not null; default:'0'"` // 是否必填题，0 必填   1 非必填
	QuestionPicID int // 图片id
}
