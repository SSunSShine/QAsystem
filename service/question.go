package service

import "github.com/SSunSShine/QAsystem/model"

// CreateQuestionInterface struct
type CreateQuestionInterface struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func (cq *CreateQuestionInterface) Create(ProfileID uint) (q model.Question, code int, err error) {

	var p model.Profile

	q.Title = cq.Title
	q.Desc = cq.Desc
	q.QuestionProfileID = ProfileID
	p.ID = ProfileID
	q.QuestionProfile, _ = p.Get()

	code, err = q.Create()

	return
}