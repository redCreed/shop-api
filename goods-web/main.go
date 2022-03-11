package main

import (
	"flag"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/initialize"
	"mxshop-api/goods-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//解析命令行参数
	var mode string
	flag.StringVar(&mode, "m", "develop", "-m config development")
	flag.Parse()
	//初始化配置
	initialize.InitConfig()
	//初始化日志
	initialize.InitLog()
	//初始化grpc对象与客户端负载均衡处理
	initialize.InitSrvConn2()
	//初始化验证器
	if err :=initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	//固定端口
	//port, err := utils.GetFreePort()
	//if err == nil {
	//	global.ServerConfig.Port = port
	//}

	//初始化路由
	engine := initialize.InitRouter()
	//服务注册
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	serviceId := uuid.String()
	//注入consul地址
	registerClient := consul.NewRegistry(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	//注册api服务地址到consul
	registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port,
		global.ServerConfig.Name, global.ServerConfig.Tags,serviceId)

	go func() {
		if err := engine.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			panic(err)
		}
	}()
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	}
	fmt.Println("注销成功")
}

