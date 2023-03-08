package global

import "fmt"

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

var HandlerErrors = struct {
	ValidateError    CustomError
	ClaimsTokenError CustomError
}{
	ValidateError:    CustomError{40001, "请求参数错误"},
	ClaimsTokenError: CustomError{40002, "登录鉴权失败"},
}

var ServiceErrors = struct {
	ServiceError CustomError
	ParamError   CustomError
}{
	ServiceError: CustomError{50001, "service错误"},
	ParamError:   CustomError{50002, "service调用参数错误"},
}

func (ce CustomError) Error() string {
	return fmt.Sprint(ce.ErrorCode, ", ", ce.ErrorMsg)
}
