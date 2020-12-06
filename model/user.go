package model

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// User 用户
type User struct {
	gorm.Model
	Mail      string  `json:"mail" gorm:"type:varchar(100);unique_index;not null"`
	Password  string  `json:"password" gorm:"type:varchar(100);not null"`
	Phone     string  `json:"phone" gorm:"type:varchar(11);not null"`
}

func (u *User) Get() (user User, err error) {

	if err = database.DB.Where(&u).First(&user).Error; err != nil {
		log.Print(err)
	}

	return
}

func (u *User) Create() (err error) {

	if err = database.DB.Create(&u).Error; err != nil {
		log.Print(err)
	}

	return
}

func (u *User) Update() (err error) {

	if err = database.DB.Model(&u).Updates(u).Error; err != nil {
		log.Print(err)
	}

	return
}

func (u *User) Delete() (err error) {

	if err = database.DB.Unscoped().Delete(&u).Error; err != nil {
		log.Print(err)
	}

	return
}

func (u *User) Count() (count int, err error) {

	if err = database.DB.Model(&u).Count(&count).Error; err != nil {
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

// AfterUpdate 密码加密
func (u *User) AfterUpdate(db *gorm.DB) (err error) {

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
