package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

// Question 问题
type Question struct {
	gorm.Model
	Title  string `json:"title" gorm:"type:varchar(50);not null"`
	Desc   string `json:"desc" gorm:"type:varchar(4000);not null"`
	UserID uint   `json:"userId"`
	User   User   `json:"user"  gorm:"ForeignKey:UserID"`
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

func (q *Question) Delete() ( err error) {

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
