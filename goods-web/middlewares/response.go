package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/global"
	myerror "mxshop-api/user-web/middlewares/error"
)

type Response struct {
	Ret  int64       `json:"ret"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func ResponseFail(ctx *gin.Context) {
	data := make([]int, 0)
	resp := &Response{Ret: -1, Msg: "failed", Data: data}
	ctx.JSON(200, resp)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	if data == "" || data == nil {
		data = make([]int, 0)
	}
	resp := &Response{Ret: 0, Msg: "成功", Data: data}
	c.JSON(200, resp)
}

func ResponseData(ctx *gin.Context, code int, msg string) {
	data := make([]int, 0)
	resp := &Response{Ret: int64(code), Msg: msg, Data: data}
	ctx.JSON(200, resp)
}

//返回响应
func MyResponseData(c *gin.Context, err error, data interface{}) {
	if data == "" || data == nil {
		data = make([]int, 0)
	}
	resp := &Response{Data: data, Ret: 0, Msg: "操作成功"}
	if err != nil { // 这里是错误模式
		if tem, ok := err.(*myerror.MyError); ok {
			resp = &Response{Ret: tem.Code, Msg: tem.Error(), Data: make([]int, 0)}
		} else {
			// 包装一下
			stack := fmt.Sprintf("stack error trace:\n%+v\n", err) //错误的堆栈
			fmt.Println(stack)
			global.ZapLog.Error(stack)
			//返回给前端
			resp = &Response{Ret: myerror.COMMENT_ERROR, Msg: myerror.ErrText[myerror.COMMENT_ERROR], Data: make([]int, 0)}
		}
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, resp)
	c.Abort()
}
