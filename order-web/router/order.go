package router

import (
	"github.com/gin-gonic/gin"

	"mxshop-api/order-web/api/order"
	"mxshop-api/order-web/api/pay"
	"mxshop-api/order-web/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("/orders").Use(middlewares.Trace())
	{
		OrderRouter.GET("", order.List)   // 订单列表
		OrderRouter.POST("/create",  order.New)  // 新建订单
		OrderRouter.GET("/:id", order.Detail)  // 订单详情
	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}
}