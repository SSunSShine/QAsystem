package api

import (
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ProfileVO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Gender    int       `json:"gender"`
	Desc      string    `json:"desc"`
	User      UserVO    `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetProfile 获取单个用户简介
func GetProfile(c *gin.Context) {

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
	util.SimpleCopyProperties(&profileVO.User, &profile.User)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data": profileVO,
	})
}

func UpdateProfile(c *gin.Context) {

	var p model.Profile

	id, _ := strconv.Atoi(c.Param("id"))
	p.ID = uint(id)

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if err := p.Update(); err != nil {
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

func GetProfilesCount(c *gin.Context) {

	var p model.Profile

	if count, err := p.Count(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
			"data": count,
		})
	}
}

// DeleteProfile
func DeleteProfile(c *gin.Context)  {

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

	if err := profile.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
		})
	}
}