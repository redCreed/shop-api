package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/api/banners"
	"mxshop-api/goods-web/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners").Use(middlewares.Trace())
	{
		BannerRouter.GET("", banners.List)          // 轮播图列表页
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), banners.Delete) // 删除轮播图
		BannerRouter.POST("",  middlewares.JWTAuth(),  banners.New)       //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), banners.Update) //修改轮播图信息
	}
}