package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

type Comment struct {
	gorm.Model
	Content  string `json:"content" gorm:"type:varchar(500)"`
	Answer   Answer `json:"answer" gorm:"ForeignKey:AnswerID"`
	AnswerID uint   `json:"answerId"`
	UserID   uint   `json:"userId"`
	User     User   `json:"user"  gorm:"ForeignKey:UserID"`
}

func (c *Comment) Get() (comment Comment, err error) {

	if err = database.DB.Where(&c).Preload("Answer").Preload("User").First(&comment).Error; err != nil {
		log.Print(err)
	}
	return
}

func (c *Comment) Create() (err error) {

	if err = database.DB.Create(&c).Error; err != nil {
		log.Print(err)
	}
	return
}

func (c *Comment) Update() (err error) {

	if err = database.DB.Model(&c).Updates(c).Error; err != nil {
		log.Print(err)
	}
	return
}

func (c *Comment) Delete() (err error) {

	if err = database.DB.Unscoped().Where(&c).First(&c).Delete(&c).Error; err != nil {
		log.Print(err)
	}
	return
}

func (c *Comment) GetList() (comments []Comment, err error) {

	if err = database.DB.Preload("Answer").Preload("User").Find(&comments, c).Error; err != nil {
		log.Print(err)
	}
	return
}

func (c *Comment) GetOrderList(order string) (comments []Comment, err error) {

	if err = database.DB.Preload("Answer").Preload("User").Order(order).Find(&comments, c).Error; err != nil {
		log.Print(err)
	}

	return
}

func (c *Comment) Count() (count int, err error) {

	if err = database.DB.Model(&c).Where(&c).Count(&count).Error; err != nil {
		log.Print(err)
	}
	return
}

// AfterDelete 回答下评论数量 - 1
func (c *Comment) AfterDelete(db *gorm.DB) (err error) {

	var a Answer
	a.ID = c.AnswerID

	if err = db.Model(&a).UpdateColumn("comments_count", gorm.Expr("comments_count - ?", 1)).Error; err != nil {
		log.Print(err)
	}

	return
}