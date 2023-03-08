package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"monaToolBox/global"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 响应成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(
		http.StatusOK, Response{
			Code:    0,
			Message: "ok",
			Data:    data,
		},
	)
}

// Fail 响应失败 ErrorCode 不为 0 表示失败
func Fail(c *gin.Context, errorCode int, msg string) {
	c.JSON(
		http.StatusOK, Response{
			Code:    errorCode,
			Message: msg,
			Data:    nil,
		},
	)
}

// FailByError 失败响应 返回自定义错误的错误码、错误信息
func FailByError(c *gin.Context, error global.CustomError, extraMsg ...string) {
	if len(extraMsg) > 0 {
		Fail(c, error.ErrorCode, fmt.Sprint(error.ErrorMsg, extraMsg))
	} else {
		Fail(c, error.ErrorCode, error.ErrorMsg)
	}
}

// ValidateFail 请求参数验证失败
func ValidateFail(c *gin.Context, msg string) {
	Fail(c, global.HandlerErrors.ValidateError.ErrorCode, msg)
}

// ServiceFail 调用service失败
func ServiceFail(c *gin.Context, extraMsg ...string) {
	if len(extraMsg) > 0 {
		Fail(c, global.ServiceErrors.ServiceError.ErrorCode, fmt.Sprint(global.ServiceErrors.ServiceError.ErrorMsg, extraMsg))
	} else {
		Fail(c, global.ServiceErrors.ServiceError.ErrorCode, global.ServiceErrors.ServiceError.ErrorMsg)
	}
}

// ClaimsTokenFail jwt鉴权失败
func ClaimsTokenFail(c *gin.Context) {
	Fail(c, global.HandlerErrors.ClaimsTokenError.ErrorCode, global.HandlerErrors.ClaimsTokenError.ErrorMsg)
}
