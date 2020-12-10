package route

import (
	"github.com/SSunSShine/QAsystem/controller/api"
	"github.com/SSunSShine/QAsystem/middleware"
	"github.com/gin-gonic/gin"
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
	}

}
