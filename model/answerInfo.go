package model

import "time"

type AnswerInfo struct {
	ID        uint `gorm:"primarykey"`
	AnswerID  int `gorm:"not null"` // 成员ID
	SurveyID int `gorm:"not null"` // 关联调查问卷ID
	QuestionID int `gorm:"not null"` // 关联问题ID
	CreatedAt time.Time
}