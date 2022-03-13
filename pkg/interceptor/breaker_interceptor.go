package interceptor

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

//限流熔断降级
func BreakerInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		hystrix.ConfigureCommand("get/user/list", hystrix.CommandConfig{
			Timeout:                10000, //执行commend命令超时时间
			MaxConcurrentRequests:  2,     //command 的最大并发量
			RequestVolumeThreshold: 2,     //请求阈值  熔断器是否打开首先要满足这个条件；这里的设置表示至少有5个请求才进行ErrorPercentThreshold错误百分比计算
			SleepWindow:            5000,     //过多长时间，熔断器再次检测是否开启。单位毫秒
			ErrorPercentThreshold:  90,    //错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
		})

		circuitBreaker, has, err := hystrix.GetCircuit("get/user/list")
		if err != nil {
			return errors.New("get circuitBreaker err")
		}
		if !has {
			//未设置熔断器
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		fmt.Println("circuitBreaker:", circuitBreaker)

		hystrix.Do("get/user/list", func() error {
			//1、根据业务码和熔断开关进行判断是否开启熔断器
			//2、去调用
			if err :=  invoker(ctx, method, req, reply, cc, opts...); err  != nil {
				return err
			}
			return nil
		},
		//降级
		func(err error) error {
			//返回特定结果
			reply = "dd"
			return nil
		},
			
		)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
