package service

import "github.com/SSunSShine/QAsystem/model"

// CreateQuestionInterface struct
type CreateQuestionInterface struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func (cq *CreateQuestionInterface) Create(UserID uint) (q model.Question, code int, err error) {

	var u model.User

	q.Title = cq.Title
	q.Desc = cq.Desc
	q.UserID = UserID
	u.ID = UserID
	q.User, _ = u.Get()

	code, err = q.Create()

	return
}