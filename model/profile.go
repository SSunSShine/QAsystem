package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

// Profile 用户简介
type Profile struct {
	gorm.Model
	Name   string `json:"name" gorm:"type:varchar(20);not null"`
	Gender int    `json:"gender" gorm:"type:int;DEFAULT:0;not null"`
	Desc   string `json:"desc"`
	UserID uint   `json:"userId"`
	User   User   `json:"user"  gorm:"ForeignKey:UserID"`
}

func (p *Profile) Get() (profile Profile, err error) {

	if err = database.DB.Where(&p).Preload("User").First(&profile).Error; err != nil {
		log.Print(err)
	}
	return
}

func (p *Profile) Create() (err error) {

	if err = database.DB.Create(&p).Error; err != nil {
		log.Print(err)
	}
	return
}

func (p *Profile) Update() (err error) {

	if err = database.DB.Model(&p).Updates(p).Error; err != nil {
		log.Print(err)
	}
	return
}

func (p *Profile) Delete() (err error) {

	if err = database.DB.Unscoped().Delete(&p).Error; err != nil {
		log.Print(err)
	}
	return
}

func (p *Profile) Count() (count int, err error) {

	if err = database.DB.Model(&p).Count(&count).Error; err != nil {
		log.Print(err)
	}

	return
}

// AfterDelete 级联删除用户信息
func (p *Profile) AfterDelete(db *gorm.DB) (err error) {

	if err = db.Where("id = ?", p.UserID).Unscoped().Delete(&User{}).Error; err != nil {
		log.Print(err)
	}

	return
}
