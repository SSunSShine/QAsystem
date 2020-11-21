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

type QuestionVO struct {
	ID                uint      `json:"id"`
	Title             string    `json:"title"`
	Desc              string    `json:"desc"`
	QuestionProfile   ProfileVO `json:"questionProfile"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// GetQuestion 获取单个问题信息
func GetQuestion(c *gin.Context)  {

	var q model.Question

	id, _ := strconv.Atoi(c.Param("id"))
	q.ID = uint(id)

	question, err := q.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	var questionVO QuestionVO
	util.SimpleCopyProperties(&questionVO, &question)
	util.SimpleCopyProperties(&questionVO.QuestionProfile, &question.QuestionProfile)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": questionVO,
	})
}

func UpdateQuestion(c *gin.Context)  {

	var q model.Question

	id, _ := strconv.Atoi(c.Param("id"))
	q.ID = uint(id)

	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if code, err := q.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": code,
		})
	}

}

func DeleteQuestion(c *gin.Context)  {

	var q model.Question

	id, _ := strconv.Atoi(c.Param("id"))
	q.ID = uint(id)

	if code, err := q.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": code,
		})
	}
}

// GetQuestionsCount 按用户简介id统计问题数量
func GetQuestionsCount(c *gin.Context)  {

	var q model.Question

	questionProfileID, _ := strconv.Atoi(c.Query("questionProfileID"))
	q.QuestionProfileID = uint(questionProfileID)

	if count, err := q.Count(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": count,
		})
	}
}

func CreateQuestion(c *gin.Context)  {

	var cq service.CreateQuestionInterface

	ProfileID, exist := c.Get("pid")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": "Not exist",
		})
		c.Abort()
		return
	}
	value, ok := ProfileID.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": "Not uint",
		})
		c.Abort()
		return
	}

	if err := c.ShouldBindJSON(&cq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	question, code, err := cq.Create(value);
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	var questionVO QuestionVO
	util.SimpleCopyProperties(&questionVO, &question)
	util.SimpleCopyProperties(&questionVO.QuestionProfile, &question.QuestionProfile)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": map[string]interface{}{
			"record": questionVO,
			"code": code,
		},
	})

}

// GetQuestions 按条件获取问题列表
func GetQuestions(c *gin.Context)  {

	var q model.Question
	//var limit int
	//var offset int
	var err error

	questionProfileID, _ := strconv.Atoi(c.Query("questionProfileID"))
	q.QuestionProfileID = uint(questionProfileID)

	//limit, err = strconv.Atoi(c.Query("limit"))
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": http.StatusInternalServerError,
	//		"message": err.Error(),
	//	})
	//	return
	//}
	//
	//offset, err = strconv.Atoi(c.Query("offset"))
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": http.StatusInternalServerError,
	//		"message": err.Error(),
	//	})
	//	return
	//}

	questions, err := q.GetList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	count := 0
	var questionsVO []QuestionVO
	for _, q := range questions {
		var questionVO QuestionVO
		util.SimpleCopyProperties(&questionVO, &q)
		util.SimpleCopyProperties(&questionVO.QuestionProfile, &q.QuestionProfile)
		questionsVO = append(questionsVO, questionVO)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": questionsVO,
		"total": count,
	})
}