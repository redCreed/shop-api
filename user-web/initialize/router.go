package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"
	"net/http"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()
	//该路由做健康检查 空接口
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	rt := engine.Group("/v1")
	//配置跨域
	rt.Use(middlewares.Cors(), middlewares.Trace())
	router.InitUserRouter(rt)
	router.InitBaseRouter(rt)
	return engine
}
