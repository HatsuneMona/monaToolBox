package types

import "monaToolBox/global"

var (
	HandlerErrors = struct {
		Conflict global.CustomError
	}{
		Conflict: global.CustomError{ErrorCode: 40101, ErrorMsg: "短链已存在"},
	}

	ServiceErrors = struct{}{}
)
