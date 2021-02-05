package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

// Question 问题
type Question struct {
	gorm.Model
	Title        string  `json:"title" gorm:"type:varchar(50);not null"`
	Desc         string  `json:"desc" gorm:"type:varchar(4000);not null"`
	UserID       uint    `json:"userId"`
	User         User    `json:"user"  gorm:"ForeignKey:UserID"`
	AnswersCount int     `json:"answersCount"`
	ViewCount    int     `json:"viewCount"`
	Hot          float64 `json:"hot" sql:"-"`
}

func (q *Question) Get() (question Question, err error) {

	if err = database.DB.Where(&q).Preload("User").First(&question).Error; err != nil {
		log.Print(err)
	}

	return
}

func (q *Question) Create() (err error) {

	if err = database.DB.Create(&q).Error; err != nil {
		log.Print(err)
	}

	return
}

func (q *Question) Update() (err error) {

	if err = database.DB.Model(&q).Updates(q).Error; err != nil {
		log.Print(err)
	}

	return
}

func (q *Question) Delete() (err error) {

	if err = database.DB.Unscoped().Delete(&q).Error; err != nil {
		log.Print(err)
	}

	return
}

func (q *Question) Count() (count int, err error) {

	if err = database.DB.Model(&q).Where(&q).Count(&count).Error; err != nil {
		log.Print(err)
	}

	return
}

// GetList 获取问题列表
func (q *Question) GetList() (questions []Question, err error) {

	if err = database.DB.Where(&q).Preload("User").Find(&questions).Error; err != nil {
		log.Print(err)
	}

	return
}

func (q *Question) GetOrderList(order string) (questions []Question, err error) {

	if err = database.DB.Preload("User").Order(order).Find(&questions, q).Error; err != nil {
		log.Print(err)
	}

	return
}

// IncrAnswersCount 回答数+1
func (q *Question) IncrAnswersCount() (err error) {

	if err := database.DB.Model(&q).UpdateColumn("answers_count", gorm.Expr("answers_count + ?", 1)).Error; err != nil {
		log.Print(err)
	}
	return
}

// IncrView 浏览数+1
func (q *Question) IncrView() (err error) {

	if err := database.DB.Model(&q).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		log.Print(err)
	}
	return
}

// AfterDelete 级联删除问题下的回答
func (q *Question) AfterDelete(db *gorm.DB) (err error) {

	var a Answer
	a.QuestionID = q.ID

	if err = db.Where(&a).Unscoped().Delete(&a).Error; err != nil {
		log.Print(err.Error() + ": delete question and question's answer")
	}
	return
}