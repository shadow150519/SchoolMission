package model

import (
	"gorm.io/gorm"
	"time"
)

// SurveyInfo 问卷主题
// 存储了问卷的主干信息，不包括具体题目和选项
type SurveyInfo struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index"`
	SurveyName string `gorm:"not null"` // 问卷名称
	SurveyDescription string `gorm:"default:null"` // 问卷描述
	StartTime time.Time `gorm:"not null"` // 问卷开始时间，默认当前时间
	EndTime time.Time `gorm:"not null"` // 问卷结束时间，若为0则表示无时间限制
	Status	int `gorm:"not null;default:1"` // 问卷状态，1 发布， 2 暂存，3 结束,  4 删除
	SurveySort int `gorm:"not null;default:1"`  // 问卷优先级，默认为1
	TopFlag int `gorm:"not null;default:2"` // 是否置顶，1置顶，2不置顶，默认为2
	CreatorID int `gorm:"not null"`  // 创建人员ID
	UpdatorID int `gorm:"not null"'` // 更新人员ID
	SurveyPicId int    				// 问卷封面图片
}

// Create 新建任务
func (s *SurveyInfo)Create(surveyName, surveyDescription string, startTime, endTime time.Time,
	status, surveySort, topFlag, creatorId, updatorId, surveyPicId int )bool {
	sql := "INSERT INTO survey_infos" +
		"(survey_name, survey_description, start_time, end_time, status, survey_sort, top_flag, creator_id, updator_id, survey_pic_id)" +
		"VALUES(?,?,?,?,?,?,?,?,?,?)"
	result := s.Raw(sql,surveyName,surveyDescription,startTime,
		endTime,status,surveySort,topFlag,creatorId,updatorId,surveyPicId)

	return result.RowsAffected > 0
}

// Update 更新任务
func (s *SurveyInfo) Update(id int64,surveyName, surveyDescription string, startTime, endTime time.Time,
	status, surveySort, topFlag, creatorId, updatorId, surveyPicId int)bool  {
	sql := "UPDATE  survey_infos" +
		"(survey_name, survey_description, start_time, end_time, status, survey_sort, top_flag, creator_id, updator_id, survey_pic_id)" +
		"VALUES(?,?,?,?,?,?,?,?,?,?)" +
		"WHERE id=?"
	result := s.Exec(sql,surveyName,surveyDescription,startTime,
		endTime,status,surveySort,topFlag,creatorId,updatorId,surveyPicId,id)
	return result.RowsAffected > 0
}

// Show 查看任务
func (s *SurveyInfo) ShowSurveysBySurveyPattern(surveyPattern string, page, limit int)([]*SurveyInfo, error){
	surveys := []*SurveyInfo{}
	sql := "SELECT `id`, `created_at`, `updated_at`, `survey_name`, " +
		"`survey_description`, `start_time`, `end_time`, `status`, " +
		"`survey_sort`, `top_flag`, `creator_id`, `updator_id` FROM survey_infos" +
		"where survey_name like ? LIMIT ?, ?"
	result := s.Raw(sql,"%"+ surveyPattern + "%", (page - 1) * 10, limit).Find(surveys)
	if result.Error == nil {
		return surveys, nil
	} else{
		return nil, result.Error
	}
}

// Delete 删除任务
func (s *SurveyInfo) Delete(id int64)bool {
	sql := "DELETE FROM survey_infos WHERE id=?"
	result := s.Raw(sql,id)
	return result.RowsAffected > 0
}