package models

import "errors"

const (
	DataAlreadyExistErrorCode = "DATA_ALREADY_EXIST"
	DataNotFoundErrorCode     = "DATA_NOT_FOUND"
	DataValidationErrorCode   = "DATA_VALIDATION_ERROR"
	InvalidJsonErrorCode      = "INVALID_JSON"
	InvalidPathParamErrorCode = "INVALID_PATH_PARAM"
	ServerErrorCode           = "SERVER_ERROR"
)

var (
	ErrEmailAlreadyExist       = errors.New("email already exist")
	ErrMobilePhoneAlreadyExist = errors.New("mobile phone already exist")
)
