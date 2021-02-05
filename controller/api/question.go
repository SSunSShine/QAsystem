package api

import (
	"context"
	"github.com/SSunSShine/QAsystem/database"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/service"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type QuestionVO struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	Desc         string     `json:"desc"`
	Questioner   Questioner `json:"questioner"`
	AnswersCount int        `json:"answersCount"`
	ViewCount	 int		`json:"viewCount"`
	Hot			 float64	`json:"hot"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type Questioner struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// GetQuestion 获取单个问题信息
func GetQuestion(c *gin.Context) {

	var q model.Question
	var p model.Profile

	qid := c.Param("id")
	id, _ := strconv.Atoi(qid)
	q.ID = uint(id)

	question, err := q.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": question",
		})
		return
	}

	p.UserID = q.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": profile",
		})
		return
	}

	// 增加浏览量记录到redis
	err = service.IncrView(qid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": questionId",
		})
		return
	}

	var questionVO QuestionVO
	util.SimpleCopyProperties(&questionVO, &question)
	util.SimpleCopyProperties(&questionVO.Questioner, &profile)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    questionVO,
	})
}

func UpdateQuestion(c *gin.Context) {

	var q model.Question

	id, _ := strconv.Atoi(c.Param("id"))
	q.ID = uint(id)

	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error()+": bind question json",
		})
		return
	}

	if err := q.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": update question",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
		})
	}

}

func DeleteQuestion(c *gin.Context) {

	var q model.Question

	id, _ := strconv.Atoi(c.Param("id"))
	q.ID = uint(id)

	if err := q.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": delete question",
		})
	} else {
		if err := database.RDB.ZRem(context.Background(), "question", q.ID).Err(); err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
		})
	}
}

// GetQuestionsCount 按用户id统计问题数量
func GetQuestionsCount(c *gin.Context) {

	var q model.Question

	userID, _ := strconv.Atoi(c.Query("userID"))
	q.UserID = uint(userID)

	if count, err := q.Count(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": questions count",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
			"data":    count,
		})
	}
}

func CreateQuestion(c *gin.Context) {

	var cq service.CreateQuestionInterface
	var p model.Profile

	UserID, exist := c.Get("uid")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not exist: userid",
		})
		c.Abort()
		return
	}
	value, ok := UserID.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not uint: userid",
		})
		c.Abort()
		return
	}

	if err := c.ShouldBindJSON(&cq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error()+": bind question json",
		})
		return
	}

	question, err := cq.Create(value)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": create question",
		})
		return
	}

	p.UserID = question.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": profile",
		})
		return
	}

	var questionVO QuestionVO
	util.SimpleCopyProperties(&questionVO, &question)
	util.SimpleCopyProperties(&questionVO.Questioner, &profile)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    questionVO,
	})

}

// GetQuestions 按条件获取问题列表
func GetQuestions(c *gin.Context) {

	var q model.Question
	var p model.Profile
	var questions []model.Question
	var err error

	userID, _ := strconv.Atoi(c.Query("userID"))
	q.UserID = uint(userID)

	order := c.Query("order")
	if order == "view_count" {
		if questions, err = q.GetOrderList("view_count desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error()+": questions",
			})
			return
		}
	} else if order == "answers_count" {
		if questions, err = q.GetOrderList("answers_count desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error()+": questions",
			})
			return
		}
	} else if order == "create_time" {
		if questions, err = q.GetOrderList("created_at desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error()+": questions",
			})
			return
		}
	} else {
		if questions, err = q.GetList(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error()+": questions",
			})
			return
		}
	}

	count := 0
	var questionsVO []QuestionVO
	for _, q := range questions {
		p.UserID = q.UserID
		profile, err := p.Get()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error()+": profile",
			})
			return
		}

		var questionVO QuestionVO
		util.SimpleCopyProperties(&questionVO, &q)
		util.SimpleCopyProperties(&questionVO.Questioner, &profile)
		questionsVO = append(questionsVO, questionVO)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    questionsVO,
		"total":   count,
	})
}

// GetTopQ 获取热榜
func GetTopQ(c *gin.Context) {

	topQ := make(map[int]interface{})
	var questionVO QuestionVO
	var p model.Profile

	for i := 1; i <= 50; i++ {
		obj, ok := service.GetTopQ().Load(strconv.Itoa(i))
		if !ok {
			log.Print("问题不足热榜数目!")
			break
		}

		question := obj.(model.Question)
		p.ID = question.UserID
		profile, err := p.Get()
		if err != nil {
			log.Print(err.Error()+": profile")
			break
		}
		util.SimpleCopyProperties(&questionVO, &question)
		util.SimpleCopyProperties(&questionVO.Questioner, &profile)
		topQ[i] = questionVO
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    topQ,
	})
}
