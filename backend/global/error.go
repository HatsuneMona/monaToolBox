package global

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

type CustomErrors struct {
	ValidateError    CustomError
	ServiceError     CustomError
	ClaimsTokenError CustomError
}

var Errors = CustomErrors{
	ValidateError:    CustomError{40001, "请求参数错误"},
	ServiceError:     CustomError{50001, "service错误"},
	ClaimsTokenError: CustomError{40002, "登录鉴权失败"},
}
