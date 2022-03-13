package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"mxshop-api/pkg/interceptor"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
	"mxshop-api/user-web/utils/otgrpc"
)

//初始化grpc连接与客户端负载均衡 切记导入上面的包grpc-consul-resolver
func InitSrvConn2() {
	//conn, err := grpc.Dial(
	////todo 注意tag tag会作为过滤条件进行过滤服务
	//	"consul://127.0.0.1:8500/whoami?wait=14s&tag=manual",
	//	grpc.WithInsecure(),
	//	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	//)
	target := fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserSrvInfo.Name)
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithChainUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),

			//interceptor.TimeoutInterceptor(1*time.Second), //超时拦截器

			//响应重试
			//grpc_retry.UnaryClientInterceptor(
			//	grpc_retry.WithMax(2),
			//	grpc_retry.WithCodes(codes.Internal), //服务器应该设置超时时间
			//	//grpc_retry.WithPerRetryTimeout(4*time.Second), //每次重试超时时间
			//	grpc_retry.WithBackoff(grpc_retry.BackoffLinearWithJitter(time.Second, 0.1)),
			//	),

			interceptor.BreakerInterceptor(),

			),
		)

	if err != nil {
		log.Fatal(err)
	}
	userSrvClient := proto.NewUserClient(conn)
	global.UserSrvClient = userSrvClient
}

//初始化conn 全局复用
func InitSrvConn()  {
	config := api.DefaultConfig()
	//consul地址
	config.Address=fmt.Sprintf(`%s:%d`, global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(config)
	if err != nil {
		global.ZapLog.Error("init srv-conn连接", zap.Error(err))
		panic(err)
	}

	//从consul中进行服务发现
	data, err :=client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		global.ZapLog.Error("filter srv-conn连接", zap.Error(err))
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
	for _, value := range data{
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == ""{
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做

	//2. 已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//3. 一个连接多个groutine共用，性能 - 连接池
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

