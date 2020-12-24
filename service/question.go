package service

import (
	"errors"
	"github.com/SSunSShine/QAsystem/database"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

// CreateQuestionInterface struct
type CreateQuestionInterface struct {
	Title string `json:"title" validate:"required,min=1,max=50" label:"问题标题"`
	Desc  string `json:"desc" validate:"required,min=1,max=4000" label:"问题描述"`
}

var ZSetKey = "question"
const maxMessageNum = 2000
var ChData = make(chan uint, maxMessageNum)

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
	if q.User, err = u.Get(); err != nil {
		return
	}

	if err = q.Create(); err != nil {
		return
	}

	// Redis创建问题记录
	database.RDB.ZAdd(ctx, ZSetKey, &redis.Z{Score: 0, Member: q.ID})

	return
}

// IncrView 更新redis,浏览量Score +1
func IncrView(qid string) (err error) {
	// 增加浏览量记录到redis
	_, err = database.RDB.ZIncrBy(ctx, ZSetKey, 1, qid).Result()
	if err != nil {
		log.Print(err)
		return
	}
	id, err := strconv.Atoi(qid)
	if err != nil {
		return
	}
	ChData <- uint(id)

	return
}

// updateViews 异步更新MySQL
func updateViews()  {
	for {
		select {
		case updateData := <-ChData:
			var q model.Question
			q.ID = updateData
			err := q.IncrView()
			if err != nil {
				log.Print(err)
			}
		}
	}
}

var isConsumerRun = false
func RunViewConsumer() {
	// Only Run one consumer.
	if !isConsumerRun {
		go updateViews()  //开启一个消费者goroutune，作用是接收redis的改动信息，更新数据库
		isConsumerRun = true
	}
}