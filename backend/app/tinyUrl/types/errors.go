package types

import "monaToolBox/global"

var ServiceErrors = struct {
	Conflict global.CustomError
}{
	Conflict: global.CustomError{ErrorCode: 50101, ErrorMsg: "短链已存在"},
}
