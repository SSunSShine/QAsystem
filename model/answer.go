package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

// Answer 问题回答
type Answer struct {
	gorm.Model
	Content         string   `json:"content" gorm:"type:varchar(4000)"`
	QuestionID      uint     `json:"questionId"`
	Question        Question `json:"question" gorm:"ForeignKey:QuestionID"`
	UserID          uint     `json:"userId"`
	User            User     `json:"user"  gorm:"ForeignKey:UserID"`
	Voters          []Voter  `json:"-"`
	SupportersCount int      `json:"supportersCount"` // 只统计点赞数
	Voted           int      `json:"voted" gorm:"-"`  // 1 赞， 0 未投票， -1 踩
	CommentsCount   int      `json:"commentsCount"`
}

func (a *Answer) Get() (answer Answer, err error) {

	if err = database.DB.Where(&a).Preload("Question").Preload("User").First(&answer).Error; err != nil {
		log.Print(err)
	}
	return
}

func (a *Answer) Create() (err error) {

	if err = database.DB.Create(&a).Error; err != nil {
		log.Print(err)
	}
	return
}

func (a *Answer) Update() (err error) {

	if err = database.DB.Model(&a).Updates(a).Error; err != nil {
		log.Print(err)
	}
	return
}

func (a *Answer) Delete() (err error) {

	if err = database.DB.Unscoped().Where(&a).First(&a).Delete(&a).Error; err != nil {
		log.Print(err)
	}
	return
}

func (a *Answer) GetList() (answers []Answer, err error) {

	if err = database.DB.Preload("Question").Preload("User").Find(&answers, a).Error; err != nil {
		log.Print(err)
	}
	return
}

func (a *Answer) GetOrderList(order string) (answers []Answer, err error) {

	if err = database.DB.Preload("Question").Preload("User").Order(order).Find(&answers, a).Error; err != nil {
		log.Print(err)
	}

	return
}

func (a *Answer) Count() (count int, err error) {

	if err = database.DB.Model(&a).Where(&a).Count(&count).Error; err != nil {
		log.Print(err)
	}
	return
}

// GetHotAnswer 得到点赞数最多的回答
func (a *Answer) GetHotAnswer() (answer Answer, err error) {

	if err = database.DB.Preload("User").Order("supporters_count desc").First(&answer, a).Error; err != nil {
		log.Print(err)
	}
	return
}

// AfterDelete 问题下回答数量 - 1
func (a *Answer) AfterDelete(db *gorm.DB) (err error) {

	var q Question
	q.ID = a.QuestionID

	if err = db.Model(&q).UpdateColumn("answers_count", gorm.Expr("answers_count - ?", 1)).Error; err != nil {
		log.Print(err)
	}

	return
}

type QIDs struct {
	QuestionId uint `json:"question_id"`
}

// GetRandomQuestionID 随机返回20个问题（无回答的不返回） ID
func GetRandomQuestionID() (question_id []QIDs) {

	database.DB.Raw("select distinct question_id from answer ORDER BY RAND() limit 20").Find(&question_id)

	return
}

func (a *Answer) GetRandomAnswer() (answer Answer, err error) {

	if err = database.DB.Where(&a).Order("RAND()").First(&answer).Error; err != nil {
		log.Print(err)
	}

	return
}

// IncrCommentsCount 评论数 +1
func (a *Answer) IncrCommentsCount() (err error) {
	if err := database.DB.Model(&a).UpdateColumn("comments_count", gorm.Expr("comments_count + ?", 1)).Error; err != nil {
		log.Print(err)
	}
	return
}
