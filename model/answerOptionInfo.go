package model

type AnswerOptionInfo struct {
	ID        uint `gorm:"primarykey"`
	AnswerID int `gorm:'not null'`
	OptionID int `gorm:"not null"`
	OptionContent string `gorm:"not null"`
	AnswerFileID int
}
