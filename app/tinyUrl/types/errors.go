package types

import "monaToolBox/global"

var (
	HandlerErrors = struct {
		Conflict global.CustomError
		NotFound global.CustomError
	}{
		Conflict: global.CustomError{ErrorCode: 40101, ErrorMsg: "短链已存在"},
		NotFound: global.CustomError{ErrorCode: 40102, ErrorMsg: "短链不存在"},
	}

	ServiceErrors = struct{}{}
)
