package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
)

// Profile 用户简介
type Profile struct {
	gorm.Model
	Name string `json:"name"`
	Gender int `json:"gender" gorm:"type:int;DEFAULT:0"`
	Desc string `json:"desc"`
}

func (p *Profile) Get() (profile Profile, err error) {

	if err = database.DB.Where(&p).First(&profile).Error; err != nil {
		log.Print(err)
	}
	return
}

func (p *Profile) Create() (code int, err error) {

	if err = database.DB.Create(&p).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}
	return
}

func (p *Profile) Update() (code int, err error) {

	if err = database.DB.Model(&p).Updates(p).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}
	return
}

func (p *Profile) Delete() (code int, err error) {

	if err = database.DB.Delete(&p).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}
	return
}

func (p *Profile) Count() (count int, err error) {

	if err = database.DB.Model(&p).Count(&count).Error; err != nil {
		log.Print(err)
	}

	return
}