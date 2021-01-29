package route

import (
	"github.com/SSunSShine/QAsystem/controller/api"
	"github.com/SSunSShine/QAsystem/middleware"
	"github.com/SSunSShine/QAsystem/service"
	"github.com/gin-gonic/gin"
	"time"
)

func InitRouter(r *gin.Engine)  {

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	auth := r.Group("/api")
	auth.Use(middleware.JwtToken())
	{
		// User
		auth.PUT("user/:id", api.UpdateUser)
		auth.DELETE("user/:id", api.DeleteUser)

		// Profile
		auth.PUT("profile/:id", api.UpdateProfile)
		auth.DELETE("profile/:id", api.DeleteProfile)

		// Question
		auth.PUT("question/:id", api.UpdateQuestion)
		auth.DELETE("question/:id", api.DeleteQuestion)
		auth.POST("question/create", api.CreateQuestion)

		// Answer
		auth.GET("answer/:id", api.GetAnswer)
		auth.POST("answer/create", api.CreateAnswer)
		auth.PUT("answer/:id", api.UpdateAnswer)
		auth.DELETE("answer/:id", api.DeleteAnswer)

		// Voter
		auth.POST("voter/:answerID", api.CreateVoter)
		auth.DELETE("voter/:answerID", api.DeleteVoter)
	}
	router := r.Group("/api")
	{
		router.POST("user/sign", api.Sign)
		router.POST("user/login", api.Login)

		// User
		router.GET("user/:id", api.GetUser)
		router.GET("users/count", api.GetUsersCount)

		// Profile
		router.GET("profile/:id", api.GetProfile)
		router.GET("profiles/count", api.GetProfilesCount)

		// Question
		router.GET("question/:id", api.GetQuestion)
		router.GET("questions/count", api.GetQuestionsCount)
		router.GET("questions/list", api.GetQuestions)
		router.GET("questions/topQ", api.GetTopQ)
		router.GET("questions/qa", api.GetQa)

		// Answer
		router.GET("answers/count", api.GetAnswersCount)
		router.GET("answers/listByQuestion", api.GetAnswersByQuestion)
		router.GET("answers/listByUser", api.GetAnswersByUser)
		router.GET("answers/listByVoter", api.GetAnswersByVoter)
	}

	// 启动goroutine异步更新数据库
	go service.UpdateAnswersCount()
	go service.UpdateViews()
	go api.UpdateSupporters()
	go service.UpdateTopQ(time.Minute*5, 50)
}
