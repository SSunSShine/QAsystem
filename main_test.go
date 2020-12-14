package main

import (
	"github.com/SSunSShine/QAsystem/conf"
	"github.com/SSunSShine/QAsystem/middleware"
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/route"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

var server *httptest.Server
var token string

func init()  {
	var u model.User

	initAll(conf.Config())

	r := gin.Default()
	route.InitRouter(r)

	// run server using httptest
	server = httptest.NewServer(r)

	u.ID = 1
	admin, _ := u.Get()
	// 生成 jwt token
	ss, _ := middleware.Gen(admin)
	token = "Bearer " + ss
}

func TestUser(t *testing.T)  {

	e := httpexpect.New(t, server.URL)

	// 注册
	e.POST("/api/user/sign").
		WithJSON(map[string]interface{}{
		"mail": "TestUser@163.com",
		"password": "123456",
		"name": "testUser",
		"gender": 1,
		"phone": "13299999999",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 登录
	e.POST("/api/user/login").
		WithJSON(map[string]interface{}{
			"mail": "TestUser@163.com",
			"password": "123456",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 修改用户
	e.PUT("/api/user/2").
		WithJSON(map[string]interface{}{
			"mail": "TestUser2@163.com",
			"password": "123456",
		}).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

}

func TestProfile(t *testing.T)  {

	e := httpexpect.New(t, server.URL)

	// 查询个人信息
	e.GET("/api/profile/1").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 修改个人信息
	e.PUT("/api/profile/2").
		WithJSON(map[string]interface{}{
			"name": "testUser333",
			"gender": 1,
		}).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 获取个人问题数量
	e.GET("/api/questions/count").WithQuery("userID", 1).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 获取个人问题列表
	e.GET("/api/questions/list").WithQuery("userID", 1).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 注销
	e.DELETE("/api/profile/2").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

}

func TestQuestion(t *testing.T)  {

	e := httpexpect.New(t, server.URL)

	// 获取问题信息
	e.GET("/api/question/1").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 发布问题
	e.POST("/api/question/create").
		WithJSON(map[string]interface{}{
			"title": "Test Question11111",
			"desc": "This is a Test question !",
		}).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 修改问题
	e.PUT("/api/question/2").
		WithJSON(map[string]interface{}{
			"title": "Test Update Question",
			"desc": "This is an Updated question !",
		}).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 获取问题数量
	e.GET("/api/questions/count").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 获取问题列表
	e.GET("/api/questions/list").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 删除问题
	e.DELETE("/api/question/2").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")
}

func TestAnswer(t *testing.T)  {

	e := httpexpect.New(t, server.URL)

	// 获取回答信息
	e.GET("/api/answer/1").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 回答问题
	e.POST("/api/answer/create").WithQuery("questionID", 1).
		WithJSON(map[string]interface{}{
			"content": "Test Answer",
		}).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 修改回答
	e.PUT("/api/answer/2").
		WithJSON(map[string]interface{}{
			"content": "Test Answer11111",
		}).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 获取某用户回答数量
	e.GET("/api/answers/count").WithQuery("userID", 1).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 获取某用户回答列表
	e.GET("/api/answers/listByUser").WithQuery("userID", 1).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 获取某问题回答列表
	e.GET("/api/answers/listByQuestion").WithQuery("questionID", 1).
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

	// 删除回答
	// 注销
	e.DELETE("/api/answer/2").
		WithHeader("Authorization", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("message").ValueEqual("message", "success")

}