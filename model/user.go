package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"log"
	"golang.org/x/crypto/bcrypt"
)

// User 用户
type User struct {
	gorm.Model
	Mail      string  `json:"mail" gorm:"type:varchar(100);unique_index"`
	Password  string  `json:"password" gorm:"type:varchar(100)"`
	ProfileID uint     `json:"profileID"`
	Profile   Profile `json:"profile" gorm:"ForeignKey:ProfileID"`
}

func (u *User) Get() (user User, err error) {

	if err = database.DB.Where(&u).Preload("Profile").First(&user).Error; err != nil {
		log.Print(err)
	}

	return
}

func (u *User) Create() (code int, err error) {

	if err = database.DB.Create(&u).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}

	return
}

func (u *User) Update() (code int, err error) {

	if err = database.DB.Model(&u).Updates(u).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}

	return
}

func (u *User) Delete() (code int, err error) {

	if err = database.DB.Unscoped().Delete(&u).Error; err != nil {
		code = -1
		log.Print(err)
	} else {
		code = 1
	}

	return
}

func (u *User) Count() (count int, err error) {

	if err = database.DB.Model(&u).Count(&count).Error; err != nil {
		log.Print(err)
	}

	return
}

// AfterDelete 级联删除用户简介信息
func (u *User) AfterDelete(db *gorm.DB) (err error) {

	if err = db.Where("id = ?", u.ProfileID).Unscoped().Delete(&Profile{}).Error; err != nil {
		log.Print(err)
	}

	return
}

// AfterCreate 密码加密
func (u *User) AfterCreate(db *gorm.DB) (err error) {

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return
	}

	if err = db.Model(&u).UpdateColumn("password", string(password)).Error; err != nil {
		log.Print(err)
	}

	return
}