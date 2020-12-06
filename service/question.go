package service

import (
	"errors"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"log"
)

// CreateQuestionInterface struct
type CreateQuestionInterface struct {
	Title string `json:"title" validate:"required,min=1,max=50" label:"问题标题"`
	Desc  string `json:"desc" validate:"required,min=1,max=4000" label:"问题描述"`
}

func (cq *CreateQuestionInterface) Create(UserID uint) (q model.Question, err error) {

	var u model.User

	msg, err := util.Validate(cq)
	if err != nil {
		log.Println(msg)
		return q, errors.New(msg)
	}

	q.Title = cq.Title
	q.Desc = cq.Desc
	q.UserID = UserID
	u.ID = UserID
	q.User, _ = u.Get()

	err = q.Create()

	return
}