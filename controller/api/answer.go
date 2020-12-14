package api

import (
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/service"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AnswerVO struct {
	ID            uint      `json:"id"`
	Content       string    `json:"content"`
	QuestionTitle string    `json:"questionTitle"`
	Answerer      Answerer  `json:"answerer"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Answerer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func GetAnswer(c *gin.Context)  {

	var a model.Answer
	var p model.Profile

	id, _ := strconv.Atoi(c.Param("id"))
	a.ID = uint(id)

	answer, err := a.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	p.UserID = a.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	var answerVO AnswerVO
	util.SimpleCopyProperties(&answerVO, &answer)
	util.SimpleCopyProperties(&answerVO.Answerer, &profile)
	answerVO.QuestionTitle = answer.Question.Title

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    answerVO,
	})
}

func UpdateAnswer(c *gin.Context)  {

	var a model.Answer

	id, _ := strconv.Atoi(c.Param("id"))
	a.ID = uint(id)

	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if err := a.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
		})
	}
}

func DeleteAnswer(c *gin.Context) {

	var a model.Answer

	id, _ := strconv.Atoi(c.Param("id"))
	a.ID = uint(id)

	if err := a.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
		})
	}
}

// GetAnswersCount 按用户id统计回答数量
func GetAnswersCount(c *gin.Context) {

	var a model.Answer

	userID, _ := strconv.Atoi(c.Query("userID"))
	a.UserID = uint(userID)

	if count, err := a.Count(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
			"data":    count,
		})
	}
}

// GetAnswersByUser 按用户获取回答列表
func GetAnswersByUser(c *gin.Context) {

	var a model.Answer
	var p model.Profile
	var err error

	userID, _ := strconv.Atoi(c.Query("userID"))
	a.UserID = uint(userID)

	answers, err := a.GetList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	p.UserID = a.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	count := 0
	var answersVO []AnswerVO
	for _, answer := range answers {
		var answerVO AnswerVO
		util.SimpleCopyProperties(&answerVO, &answer)
		util.SimpleCopyProperties(&answerVO.Answerer, &profile)
		answerVO.QuestionTitle = answer.Question.Title

		answersVO = append(answersVO, answerVO)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    answersVO,
		"total":   count,
	})
}

// GetAnswersByQuestion 按问题获取回答列表
func GetAnswersByQuestion(c *gin.Context) {

	var a model.Answer
	var p model.Profile
	var err error

	questionID, _ := strconv.Atoi(c.Query("questionID"))
	a.QuestionID = uint(questionID)

	answers, err := a.GetList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	count := 0
	var answersVO []AnswerVO
	for _, answer := range answers {
		p.UserID = answer.UserID
		profile, err := p.Get()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}

		var answerVO AnswerVO
		util.SimpleCopyProperties(&answerVO, &answer)
		util.SimpleCopyProperties(&answerVO.Answerer, &profile)
		answerVO.QuestionTitle = answer.Question.Title

		answersVO = append(answersVO, answerVO)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    answersVO,
		"total":   count,
	})
}

func CreateAnswer(c *gin.Context) {

	var ca service.CreateAnswerInterface
	var p model.Profile

	questionID, _ := strconv.Atoi(c.Query("questionID"))
	qid := uint(questionID)

	UserID, exist := c.Get("uid")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not exist",
		})
		c.Abort()
		return
	}
	uid, ok := UserID.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not uint",
		})
		c.Abort()
		return
	}

	if err := c.ShouldBindJSON(&ca); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	answer, err := ca.Create(uid, qid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	p.UserID = answer.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	var answerVO AnswerVO
	util.SimpleCopyProperties(&answerVO, &answer)
	util.SimpleCopyProperties(&answerVO.Answerer, &profile)
	answerVO.QuestionTitle = answer.Question.Title

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    answerVO,
	})
}