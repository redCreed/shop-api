package myerror

type MyError struct {
	err  string
	Code int64
}

func (m *MyError) Error() string {
	return m.err
}

func ErrWithArgsMsg(msg string) error {
	return &MyError{
		err:  msg,
		Code: ARGS_MISS_OR_ILLEGAL_ERROR,
	}
}

func ErrWithCode(code int64) error {
	if v, ok := ErrText[code]; ok {
		return &MyError{
			err:  v,
			Code: code,
		}
	}

	return &MyError{
		err:  ErrText[COMMENT_ERROR],
		Code: COMMENT_ERROR,
	}
}

func New(str string, code int64) error {
	return &MyError{
		err:  str,
		Code: code,
	}
}

/*
	错误码
*/
const (
	FUNCTION_NOT_AVAILABLE        = 109
	COMMENT_ERROR                 = 110
	ARGS_MISS_ERROR               = 111
	ARGS_ILLEGAL_ERROR            = 112
	UPLOAD_FAIL_ERROR             = 113
	FETCH_OSS_STS_ERROR           = 114
	MISSING_WEB_TOKEN_ERROR       = 115
	VERIFY_SIGN_FAIL_ERROR        = 116
	NETWORK_REQUEST_TIMEOUT_ERROR = 117
	TIME_EXPIRED_ERROR            = 118
	FORCED_TO_LOGOFF_ERROR        = 119
	SENSITIVE_WORDS_ERROR         = 120
	NO_LESS_THAN_FIVE             = 121
	MISSING_DEVICE_INFORMATION    = 122
	USER_NOT_REGISTERED_ERROR     = 123
	USER_LOGGED_OFF_ERROR         = 124
	USER_ONE_KEY_LOGIN_FAIL       = 125
	NO_LESS_THAN_THREE            = 130
	CANNOT_BE_NEGATIVE            = 131
	ARGS_MISS_OR_ILLEGAL_ERROR    = 132
)

var ErrText = map[int64]string{
	FUNCTION_NOT_AVAILABLE:        "功能暂未开发",
	COMMENT_ERROR:                 "服务器繁忙，稍后再试！",
	ARGS_MISS_ERROR:               "参数缺失！",
	ARGS_ILLEGAL_ERROR:            "参数非法！",
	ARGS_MISS_OR_ILLEGAL_ERROR:    "参数非法或缺失！",
	UPLOAD_FAIL_ERROR:             "文件上传失败！",
	FETCH_OSS_STS_ERROR:           "获取sts授权失败！",
	MISSING_WEB_TOKEN_ERROR:       "缺失参数webToken！",
	VERIFY_SIGN_FAIL_ERROR:        "验证签名失败！",
	NETWORK_REQUEST_TIMEOUT_ERROR: "网络请求超时！",
	TIME_EXPIRED_ERROR:            "账号已过期！",
	FORCED_TO_LOGOFF_ERROR:        "此账号已在其它设备登录！",
	SENSITIVE_WORDS_ERROR:         "内容包含敏感词！",
	NO_LESS_THAN_FIVE:             "不少于5！",
	NO_LESS_THAN_THREE:            "擅长标签不能为空且最多三个！",
	CANNOT_BE_NEGATIVE:            "不能为负数！",
	MISSING_DEVICE_INFORMATION:    "缺少设备信息！",
	USER_ONE_KEY_LOGIN_FAIL:       "一键登录授权失败！",
	USER_NOT_REGISTERED_ERROR:     "用户未注册",
	USER_LOGGED_OFF_ERROR:         "当前用户已注销",
}
