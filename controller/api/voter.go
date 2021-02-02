package api

import (
	"context"
	"github.com/SSunSShine/QAsystem/database"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const maxMessageNum = 2000
var answerChan = make(chan uint, maxMessageNum)

func CreateVoter(c *gin.Context)  {

	var v model.Voter
	var a model.Answer
	var err error

	UserID, exist := c.Get("uid")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not exist: userid",
		})
		c.Abort()
		return
	}
	uid, ok := UserID.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not uint: userid",
		})
		c.Abort()
		return
	}
	v.UserID = uid

	id, _ := strconv.Atoi(c.Param("answerID"))
	v.AnswerID = uint(id)
	v.UpOrDown, err = strconv.ParseBool(c.Query("upOrDown"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if err = service.AddVoter(v.AnswerID, v.UserID, v.UpOrDown); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": add Voter",
		})
		return
	}

	if v.UpOrDown {
		// 通知异步更新 MySQL answer表
		answerChan <- v.AnswerID

		a.ID = v.AnswerID
		answer, err := a.Get()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusNotFound,
				"message": err.Error()+": answer",
			})
			return
		}

		// 增加热度记录到redis 点赞*2
		_, err = database.RDB.ZIncrBy(context.Background(), service.ZSetKey, 2, strconv.Itoa(int(answer.QuestionID))).Result()
		if err != nil {
			log.Print(err)
			return
		}
	}

	if err = v.Create(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": create voter",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
		})
	}
	return
}

func DeleteVoter(c *gin.Context)  {

	var v model.Voter
	var err error

	UserID, exist := c.Get("uid")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not exist: userid",
		})
		c.Abort()
		return
	}
	uid, ok := UserID.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not uint: userid",
		})
		c.Abort()
		return
	}
	v.UserID = uid

	id, _ := strconv.Atoi(c.Param("answerID"))
	v.AnswerID = uint(id)
	v.UpOrDown, err = strconv.ParseBool(c.Query("upOrDown"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if err = service.RemoveVoter(v.AnswerID, v.UserID, v.UpOrDown); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": remove voter",
		})
		return
	}

	if err = v.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": delete voter",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
		})
	}
	return
}

// UpdateSupporters 异步更新MySQL
func UpdateSupporters()  {
	for {
		select {
		case updateData := <-answerChan:
			var a model.Answer
			a.ID = updateData
			err := a.IncrSupporters()
			if err != nil {
				log.Print(err)
			}
		}
	}
}