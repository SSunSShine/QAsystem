package service

import "github.com/SSunSShine/QAsystem/model"

// SignInterface 注册接口
type SignInterface struct {
	Name     string `json:"name"`
	Gender	 int 	`json:"gender"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func (s *SignInterface) Sign() (code int, err error) {

	var p model.Profile
	var u model.User

	u.Mail = s.Mail
	u.Password = s.Password
	code, err = u.Create()

	p.Name = s.Name
	p.Desc = "Hello ~ " + p.Name
	p.Gender = s.Gender
	p.UserID = u.ID
	code, err = p.Create()

	return
}
