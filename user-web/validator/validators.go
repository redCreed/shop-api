package validator

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mxshop-api/user-web/global"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok{
		return false
	}
	return true
}


func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}


//校验post请求参数
func CheckPostParams(c *gin.Context, params interface{}) error {
	var (
		bts []byte
		err error
	)

	if err = c.ShouldBindJSON(params); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			return  errs
		}
		// validator.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		resultMap := removeTopStruct(errs.Translate(global.Trans))
		if bts,err =json.Marshal(resultMap); err != nil {
			return err
		}
		return errors.New(string(bts))
	}

	return nil
}