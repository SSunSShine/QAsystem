package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

// Answer 问题回答
type Answer struct {
	gorm.Model
	Content    string   `json:"content" gorm:"type:varchar(4000)"`
	QuestionID uint     `json:"questionId"`
	Question   Question `json:"question" gorm:"ForeignKey:QuestionID"`
	UserID     uint     `json:"userId"`
	User       User     `json:"user"  gorm:"ForeignKey:UserID"`
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

func (a *Answer) Count() (count int, err error) {

	if err = database.DB.Model(&a).Where(&a).Count(&count).Error; err != nil {
		log.Print(err)
	}
	return
}

// AfterCreate 问题下回答数量 + 1
func (a *Answer) AfterCreate(db *gorm.DB) (err error) {

	var q Question
	q.ID = a.QuestionID

	if err = db.Model(&q).UpdateColumn("answers_count", gorm.Expr("answers_count + ?", 1)).Error; err != nil {
		log.Print(err)
	}
	return
}

// AfterCreate 问题下回答数量 - 1
func (a *Answer) AfterDelete(db *gorm.DB) (err error) {

	var q Question
	q.ID = a.QuestionID

	if err = db.Model(&q).UpdateColumn("answers_count", gorm.Expr("answers_count - ?", 1)).Error; err != nil {
		log.Print(err)
	}
	return
}
