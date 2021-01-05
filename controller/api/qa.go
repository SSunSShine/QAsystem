package api

import (
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type QA struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	Desc         string     `json:"desc"`
	Questioner   Questioner `json:"questioner"`
	AnswersCount int        `json:"answersCount"`
	ViewCount	 int		`json:"viewCount"`
	HotAnswer	 AnswerVO	`json:"hotAnswer"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// GetQA 获取带最热回答的所有问题列表，注：不返回没有回答的问题
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
		util.SimpleCopyProperties(&qa.HotAnswer, &answer)
		util.SimpleCopyProperties(&qa.HotAnswer.Answerer, &answerProfile)

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

