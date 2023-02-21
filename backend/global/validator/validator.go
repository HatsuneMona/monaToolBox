package validator

import "github.com/go-playground/validator/v10"

// Validator 实现这个接口，即可输出自定义信息
type Validator interface {
	GetMessages() ValidatorMassages
}

// ValidatorMassages 自定义错误信息
// map[Field.Tag] = errorMsg
type ValidatorMassages map[string]string

// GetErrorMsg 获取错误信息
func GetErrorMsg(request interface{}, err error) string {

	// 断言，如果error是 validator.ValidationErrors
	if _, isValidatorErrors := err.(validator.ValidationErrors); isValidatorErrors {

		// 再判断 request 是否 实现了 Validator 接口
		_, isValidator := request.(Validator)

		for _, v := range err.(validator.ValidationErrors) {
			// 若 request 结构体实现 Validator 接口 则输出自定义错误信息
			if isValidator {
				if message, exist := request.(Validator).GetMessages()[v.Field()+"."+v.Tag()]; exist {
					return message
				}
			}
			return v.Error()
		}
	}
	return err.Error()
}
