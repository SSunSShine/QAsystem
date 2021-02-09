package api

import (
	"github.com/SSunSShine/QAsystem/middleware"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/service"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserVO struct {
	ID        uint      `json:"id"`
	Mail      string    `json:"mail"`
	Password  string    `json:"password"`
	Phone     string	`json:"phone"`
}

// GetUser 获取单个用户信息
func GetUser(c *gin.Context) {

	var u model.User

	id, _ := strconv.Atoi(c.Param("id"))
	u.ID = uint(id)

	user, err := u.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": user",
		})
		return
	}
	var userVO UserVO
	err = util.SimpleCopyProperties(&userVO, &user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error()+": copy userVO",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "success",
		"data": userVO,
	})

}

// UpdateUser 更新用户信息。 注：更新Profile个人简介需要调用UpdateProfile
func UpdateUser(c *gin.Context)  {

	var u model.User

	id, _ := strconv.Atoi(c.Param("id"))
	u.ID = uint(id)

	user, err := u.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error()+": userid",
		})
		return
	}
	uid, _ := c.Get("uid")
	UID := uid.(uint)
	if UID != user.ID {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "无权修改他人的信息",
		})
		return
	}

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": bind user json",
		})
		return
	}
	// 防止json中的id 与 url的id不同
	if user.ID != u.ID {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "JSON中的id与url中的id不同",
		})
		return
	}

	if err := u.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error()+": update user",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
		})
	}
}

// DeleteUser
func DeleteUser(c *gin.Context)  {
	
	var u model.User

	id, _ := strconv.Atoi(c.Param("id"))
	u.ID = uint(id)

	user, err := u.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error()+": user",
		})
		return
	}
	uid, _ := c.Get("uid")
	UID := uid.(uint)
	if UID != user.ID {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "无权修改他人的信息",
		})
		return
	}

	if err := user.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error()+": delete user",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
		})
	}
}

// GetUsersCount
func GetUsersCount(c *gin.Context)  {

	var u model.User

	if count, err := u.Count(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error()+": users count",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
			"data": count,
		})
	}
}

// Sign 用户注册
func Sign(c *gin.Context)  {

	var s service.SignInterface

	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": bind user json",
		})
		return
	}

	if err := s.Sign(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": sign",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "success",
		})
	}
}

// Login 用户登录
func Login(c *gin.Context)  {

	var l service.LoginInterface
	var p model.Profile

	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": bind user json",
		})
		return
	}

	user, err := l.Login()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": login",
		})
		return
	}

	p.UserID = user.ID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"message": err.Error()+": profile",
		})
		return
	}

	token, err := middleware.Gen(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"message": err.Error()+": generate token",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "success",
		"token": token,
		"uid": user.ID,
		"pid": profile.ID,
	})
}