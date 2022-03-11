package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/list", api.GetUserList)
		userRouter.POST("/pwd_login", api.PassWordLogin)
		userRouter.POST("register", api.Register)

		userRouter.GET("detail", middlewares.JWTAuth(), api.GetUserDetail)
		userRouter.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)
	}
}
