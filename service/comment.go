package service

import (
	"errors"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"log"
)

// CreateCommentInterface struct
type CreateCommentInterface struct {
	Content    string `json:"content" validate:"required,min=1,max=500" label:"回答内容"`
}

var commentsCountChan = make(chan uint, 20000)

func (cc *CreateCommentInterface) Create(UserID, AnswerID uint) (co model.Comment, err error) {

	var u model.User
	var a model.Answer

	msg, err := util.Validate(cc)
	if err != nil {
		log.Println(msg)
		return co, errors.New(msg)
	}

	co.Content = cc.Content
	co.AnswerID = AnswerID
	a.ID = AnswerID
	co.Answer, _ = a.Get()

	co.UserID = UserID
	u.ID = UserID
	co.User, _ = u.Get()

	err = co.Create()
	if err != nil {
		return
	}

	commentsCountChan <- AnswerID

	return
}

// UpdateCommentsCount 异步更新MySQL
func UpdateCommentsCount()  {
	for {
		select {
		case updateData := <-commentsCountChan:
			var a model.Answer
			a.ID = updateData
			err := a.IncrCommentsCount()
			if err != nil {
				log.Print(err)
			}
		}
	}
}