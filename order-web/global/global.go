package global

import (
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"
	"mxshop-api/order-web/config"
	"mxshop-api/order-web/proto"
)

var (
	//全局日志对象
	ZapLog *zap.Logger
	//配置对线
	//Config *models.Config
	//参数校验器
	Trans         ut.Translator
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}

	GoodsSrvClient proto.GoodsClient

	OrderSrvClient proto.OrderClient

	InventorySrvClient proto.InventoryClient
)
