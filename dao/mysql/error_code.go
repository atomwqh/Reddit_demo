package mysql

import "errors"

var (
	ErrorUserExist     = errors.New("user already exist")
	ErrorUserNotExist  = errors.New("user not exist")
	ErrorWrongPassword = errors.New("uid or password is wrong")
	ErrorInvaildID     = errors.New("invaild uid")
)
