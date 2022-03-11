package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/middlewares"
	myerror "mxshop-api/user-web/middlewares/error"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"mxshop-api/user-web/validator"
	"net/http"
	"reflect"
	"time"
)

func HandleGrpcErrorToHttp(ctx *gin.Context, err error) {
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{"msg": status.Message()})
			case codes.Internal:
				ctx.JSON(http.StatusNotFound, gin.H{"msg": "内部错误"})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "其他错误"})
			}
		}
	}
}

func GetUserList(ctx *gin.Context) {
	page := new(models.Page)
	if err := ctx.ShouldBindQuery(page); err != nil {
		middlewares.ResponseData(ctx, -1, "参数非法")
		return
	}
	tracerCtx := context.WithValue(context.Background(), "ginContext", ctx)
	resp, err := global.UserSrvClient.GetUserList(tracerCtx, &proto.PageInfo{
		Pn:   uint32(page.Page) ,
		PSize: uint32(page.PageSize) ,
	})

	if err != nil {
		global.ZapLog.Error("获取用户列表数据", zap.Error(err))
		HandleGrpcErrorToHttp(ctx, err)
		return
	}

	data := new(models.UserInfoList)
	data.Total = resp.Total
	for _, v := range resp.Data {
		temp := models.UserInfo{
			Id:       v.Id,
			Birthday: v.BirthDay,
			Mobile:   v.Mobile,
			NickName: v.NickName,
			Password: v.PassWord,
			Role:     v.Role,
		}
		data.UserInfo = append(data.UserInfo, temp)
	}

	ctx.JSON(200, data)
}


func PassWordLogin(c *gin.Context){
	//表单验证
	passwordLoginForm := &forms.PassWordLoginForm{}
	if err := validator.CheckPostParams(c, passwordLoginForm); err != nil {
		fmt.Println( reflect.TypeOf(myerror.ErrWithArgsMsg(err.Error()) ) )
		middlewares.MyResponseData(c,myerror.ErrWithArgsMsg(err.Error()), "")
		return
	}

	//if store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false){
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"captcha":"验证码错误",
	//	})
	//	return
	//}

	//登录的逻辑
	if rsp, err :=  global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile":"用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile":"登录失败",
				})
			}
			return
		}
	}else{
		//只是查询到用户了而已，并没有检查密码
		if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password: passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password":"登录失败",
			})
		}else{
			if passRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:             uint(rsp.Id),
					NickName:       rsp.NickName,
					AuthorityId:    uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(), //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer: "imooc",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":"生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id": rsp.Id,
					"nick_name": rsp.NickName,
					"token": token,
					"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
				})
			}else{
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg":"登录失败",
				})
			}
		}
	}
}


func Register(c *gin.Context){
	//用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		middlewares.MyResponseData(c,myerror.ErrWithArgsMsg(err.Error()), "")
		return
	}

	//验证码 暂时不要了
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	//})
	//value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	//if err == redis.Nil{
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"code":"验证码错误",
	//	})
	//	return
	//}else{
	//	if value != registerForm.Code {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"code":"验证码错误",
	//		})
	//		return
	//	}
	//}

	user, err :=  global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:             uint(user.Id),
		NickName:       user.NickName,
		AuthorityId:    uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer: "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":"生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.Id,
		"nick_name": user.NickName,
		"token": token,
		"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
	})
}

func GetUserDetail(ctx *gin.Context){
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)

	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		HandleGrpcErrorToHttp(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"name":rsp.NickName,
		"birthday": time.Unix(int64(rsp.BirthDay), 0).Format("2006-01-02"),
		"gender":rsp.Gender,
		"mobile":rsp.Mobile,
	})
}


func UpdateUser(ctx *gin.Context){
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		middlewares.MyResponseData(ctx,myerror.ErrWithArgsMsg(err.Error()), "")
		return
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	//将前端传递过来的日期格式转换成int
	loc, _ := time.LoadLocation("Local") //local的L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc)

	_, err :=  global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
		NickName: updateUserForm.Name,
		Gender:   updateUserForm.Gender,
		BirthDay: uint64(birthDay.Unix()),
	})
	if err != nil {
		HandleGrpcErrorToHttp(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}