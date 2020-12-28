package api

import (
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
	var err error

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
			"message": err.Error(),
		})
		return
	}

	if v.UpOrDown {
		// 通知异步更新 MySQL answer表
		answerChan <- v.AnswerID
	}

	if err = v.Create(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error(),
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
			"message": err.Error(),
		})
		return
	}

	if err = v.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error(),
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