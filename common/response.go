package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	myerror "mxshop-api/user-web/middlewares/error"
)

//ret 0 正常  msg   data
type Response struct {
	Ret  int64       `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}



//返回响应
func ResponseData(c *gin.Context, err error, data interface{}) {
	if data == "" || data == nil {
		data = make([]int, 0)
	}
	resp := &Response{Data: data, Ret: 0, Msg: "操作成功"}
	if err != nil { // 这里是错误模式
		if tem, ok := err.(*myerror.MyError); ok {
			resp = &Response{Ret: tem.Code, Msg: myerror.ErrText[tem.Code], Data: make([]int, 0)}
		} else {
			// 包装一下
			stack := fmt.Sprintf("stack error trace:\n%+v\n", err) //错误的堆栈
			fmt.Println(stack)
			public.ZapLog.Error(stack)
			//返回给前端
			resp = &Response{Ret: myerror.COMMENT_ERROR, Msg: myerror.ErrText[myerror.COMMENT_ERROR], Data: make([]int, 0)}
		}
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, resp)
	c.Abort()
}
