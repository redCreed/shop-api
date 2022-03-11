package global

import (
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"
	"mxshop-api/goods-web/config"
	"mxshop-api/goods-web/proto"
)

var (
	//全局日志对象
	ZapLog *zap.Logger
	//配置对线
	//Config *models.Config
	//参数校验器
	Trans         ut.Translator
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
)
