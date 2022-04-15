package model

import "errors"

// 根据业务逻辑需要，自定义一些错误
var (
	ERROR_USER_DOES_NOT_EXISTS = errors.New("User does not exists!!")
	ERROR_USER_PASSWORD        = errors.New("Password is invalid,Please check your password!")

	// register
	ERROR_USER_ALREADY_EXISTS     = errors.New("UserName already exists!")
	ERROR_PASSWORD_DOES_NOT_MATCH = errors.New("Password is not match!")
)
