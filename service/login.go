package service

import (
	"github.com/SSunSShine/QAsystem/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// LoginInterface 登录接口
type LoginInterface struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func (l *LoginInterface) Login() (code int, user model.User, err error) {

	var u model.User

	u.Mail = l.Mail

	user, err = u.Get()
	if err != nil {
		code = -1
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.Password)); err != nil {
		code = -1
		log.Print(err)
		return
	}

	code = 1
	return
}

