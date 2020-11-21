package api

import (
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProfileVO struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Gender int    `json:"gender"`
	Desc   string `json:"desc"`
}


// GetProfile 获取单个用户简介
func GetProfile(c *gin.Context)  {

	var p model.Profile

	id, _ := strconv.Atoi(c.Param("id"))
	p.ID = uint(id)

	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}
	var profileVO ProfileVO
	util.SimpleCopyProperties(&profileVO, &profile)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": profileVO,
	})
}

func UpdateProfile(c *gin.Context)  {

	var p model.Profile

	id, _ := strconv.Atoi(c.Param("id"))
	p.ID = uint(id)

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if code, err := p.Update(); err != nil {
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

func GetProfilesCount(c *gin.Context)  {

	var p model.Profile

	if count, err := p.Count(); err != nil {
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