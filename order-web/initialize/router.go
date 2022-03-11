package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/order-web/router"
	"mxshop-api/order-web/middlewares"
	"net/http"
)

func InitRouter() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())
	//添加链路追踪
	ApiGroup := Router.Group("/g/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)

	return Router
}
