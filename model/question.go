package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

// Question 问题
type Question struct {
	gorm.Model
	Title             string  `json:"title"`
	Desc              string  `json:"desc" gorm:"type:varchar(4000)"`
	QuestionProfileID uint    `json:"questionProfileID"`
	QuestionProfile   Profile `json:"questionProfile" gorm:"ForeignKey:QuestionProfileID"`
}

func (q *Question) Get() (question Question, err error) {

	if err = database.DB.Where(&q).Preload("QuestionProfile").First(&question).Error; err != nil {
		log.Print(err)
	}

	return
}

func (q *Question) Create() (code int, err error) {

	if err = database.DB.Create(&q).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}

	return
}

func (q *Question) Update() (code int, err error) {

	if err = database.DB.Model(&q).Updates(q).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}

	return
}

func (q *Question) Delete() (code int, err error) {

	if err = database.DB.Unscoped().Delete(&q).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
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

	if err = database.DB.Where(&q).Preload("QuestionProfile").Find(&questions).Error; err != nil {
		log.Print(err)
	}

	return
}
