package api

import (
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// QA 首页问答 （目前暂时带最热的回答）
type QA struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	Desc         string     `json:"desc"`
	Questioner   Questioner `json:"questioner"`
	AnswersCount int        `json:"answersCount"`
	ViewCount    int        `json:"viewCount"`
	Answer       AnswerVO   `json:"answer"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// 暂时弃用
// GetQA 获取带最热回答的所有问题列表，注：不返回没有回答的问题 （可做首页推荐）
func GetQA(c *gin.Context) {

	var q model.Question
	var p model.Profile
	var a model.Answer
	var questions []model.Question
	var err error

	userID, _ := strconv.Atoi(c.Query("userID"))
	q.UserID = uint(userID)

	order := c.Query("order")
	if order == "view_count" {
		if questions, err = q.GetOrderList("view_count desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}
	} else if order == "answers_count" {
		if questions, err = q.GetOrderList("answers_count desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}
	} else if order == "create_time" {
		if questions, err = q.GetOrderList("created_at desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}
	} else {
		if questions, err = q.GetList(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}
	}

	count := 0
	var qas []QA
	for _, q := range questions {
		if q.AnswersCount == 0 {
			continue
		}
		p.UserID = q.UserID
		questionProfile, err := p.Get()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}

		// 封装问题
		var qa QA
		util.SimpleCopyProperties(&qa, &q)
		util.SimpleCopyProperties(&qa.Questioner, &questionProfile)

		// 封装最热回答(点赞最多的)
		a.QuestionID = qa.ID
		answer, err := a.GetHotAnswer()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}

		p.UserID = answer.UserID
		answerProfile, err := p.Get()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}
		util.SimpleCopyProperties(&qa.Answer, &answer)
		util.SimpleCopyProperties(&qa.Answer.Answerer, &answerProfile)

		qas = append(qas, qa)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    qas,
		"total":   count,
	})
}

// GetQa 获取带随机回答的随机问题列表，注：不返回没有回答的问题 （可做首页推荐）
func GetQa(c *gin.Context) {

	var q model.Question
	var p model.Profile
	var a model.Answer

	// 随机获取 20 条有回答的问题id (不重复)
	questionIds := model.GetRandomQuestionID()

	count := 0
	var qas []QA
	for _, qid := range questionIds {
		q.ID = qid.QuestionId
		question, err := q.Get()
		if err != nil {
			log.Print("QA模块: questionId出错...")
			return
		}
		p.UserID = question.UserID
		questionProfile, err := p.Get()
		if err != nil {
			log.Print("QA模块: questionProfile出错...")
			return
		}

		// 封装问题
		var qa QA
		util.SimpleCopyProperties(&qa, &question)
		util.SimpleCopyProperties(&qa.Questioner, &questionProfile)

		// 封装随机回答
		a.QuestionID = qa.ID
		answer, err := a.GetRandomAnswer()
		if err != nil {
			log.Print("QA模块: 随机回答出错...")
			return
		}

		p.UserID = answer.UserID
		answerProfile, err := p.Get()
		if err != nil {
			log.Print("QA模块: answerProfile出错...")
			return
		}
		util.SimpleCopyProperties(&qa.Answer, &answer)
		util.SimpleCopyProperties(&qa.Answer.Answerer, &answerProfile)

		qas = append(qas, qa)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    qas,
		"total":   count,
	})
}
