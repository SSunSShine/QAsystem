package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

type Voter struct {
	gorm.Model
	AnswerID uint   `json:"answerId"`
	Answer   Answer `json:"-" gorm:"ForeignKey:AnswerID"`
	UserID   uint   `json:"userId"`
	User     User   `json:"user" gorm:"ForeignKey:UserID"`
	UpOrDown bool	`json:"upOrDown"`  // true 赞， false 踩
}

func (v *Voter) Create() (err error) {

	if err = database.DB.Create(&v).Error; err != nil {
		log.Print(err)
	}
	return
}

func (v *Voter) Delete() (err error) {

	if err = database.DB.Unscoped().Delete(&v).Error; err != nil {
		log.Print(err)
	}
	return
}

// GetList 点赞列表
func (v *Voter) GetList() (voters []Voter, err error) {

	v.UpOrDown = true
	if err = database.DB.Where(&v).Preload("Answer").Preload("Answer.Question").Find(&voters).Error; err != nil {
		log.Print(err)
	}
	return
}

// IncrSupporters 点赞数+1
func (a *Answer) IncrSupporters() (err error) {

	if err := database.DB.Model(&a).UpdateColumn("supporters_count", gorm.Expr("supporters_count + ?", 1)).Error; err != nil {
		log.Print(err)
	}
	return
}

func (v *Voter) AfterDelete(db *gorm.DB) (err error) {

	var a Answer
	a.ID = v.AnswerID
	if v.UpOrDown {
		if err = db.Model(&a).UpdateColumn("supporters_count", gorm.Expr("supporters_count - ?", 1)).Error; err != nil {
			log.Print(err)
		}
	}

	return
}
