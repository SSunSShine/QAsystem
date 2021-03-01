package service

import (
	"errors"
	"github.com/SSunSShine/QAsystem/database"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"log"
	"strconv"
)

// CreateAnswerInterface struct
type CreateAnswerInterface struct {
	Content    string `json:"content" validate:"required,min=1,max=4000" label:"回答内容"`
}

var answerCountChan = make(chan uint, 2000)

func (ca *CreateAnswerInterface) Create(UserID, QuestionID uint) (a model.Answer, err error) {

	var u model.User
	var q model.Question

	msg, err := util.Validate(ca)
	if err != nil {
		log.Println(msg)
		return a, errors.New(msg)
	}

	a.Content = ca.Content
	a.QuestionID = QuestionID
	q.ID = QuestionID
	a.Question, err = q.Get()
	if err != nil {
		log.Print(err)
		return
	}

	a.UserID = UserID
	u.ID = UserID
	a.User, _ = u.Get()

	err = a.Create()
	if err != nil {
		return
	}

	// 增加热度记录到redis 回答*0.5*1000
	_, err = database.RDB.ZIncrBy(ctx, ZSetKey, 500, strconv.Itoa(int(q.ID))).Result()
	if err != nil {
		log.Print(err)
		return
	}
	answerCountChan <- QuestionID

	return
}

// UpdateAnswersCount 异步更新MySQL
func UpdateAnswersCount()  {
	for {
		select {
		case updateData := <-answerCountChan:
			var q model.Question
			q.ID = updateData
			err := q.IncrAnswersCount()
			if err != nil {
				log.Print(err)
			}
		}
	}
}

