package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

//超时客户端拦截器
//grpc如何实现超时的
//gRPC 框架确实是通过 HTTP2 Header  中的 grpc-timeout 字段来实现跨进程传递超时时间
//总结:
//1客户端客户端发起 RPC 调用时传入了带 timeout 的 ctx
//2gRPC 框架底层通过 HTTP2 协议发送 RPC 请求时，将 timeout 值写入到 grpc-timeout HEADERS Frame 中
//3服务端接收 RPC 请求时，gRPC 框架底层解析 HTTP2 HEADERS 帧，读取 grpc-timeout 值，并覆盖透传到实际处理 RPC 请求的业务 gPRC Handle 中
//4如果此时服务端又发起对其他 gRPC 服务的调用，且使用的是透传的 ctx，这个 timeout 会减去在本进程中耗时，从而导致这个 timeout 传递到下一个 gRPC 服务端时变短，这样即实现了所谓的 超时传递 。


func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}




