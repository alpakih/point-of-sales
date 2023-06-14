package constant

import "errors"

var (
	ErrEmailAlreadyExist       = errors.New("email already exist")
	ErrMobilePhoneAlreadyExist = errors.New("mobile phone already exist")
)
