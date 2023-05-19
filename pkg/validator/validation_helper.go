package validator

import (
	"encoding/json"
	"errors"
	"github.com/beego/i18n"
	validatorGo "github.com/go-playground/validator/v10"
	"io"
)

type ErrorValidation struct {
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
	Tag   string `json:"tag,omitempty"`
	Error string `json:"error,omitempty"`
}

func HandleValidationErrors(language string, err error) interface{} {
	var validationErrors = make(map[string]interface{})

	if fields, ok := err.(validatorGo.ValidationErrors); ok {
		lang := "id"
		if language != "" {
			lang = language
		}

		if trans, found := Validate.GetTranslator(lang); found {
			for i := range fields {
				validationErrors[fields[i].Field()] = fields[i].Translate(trans)
			}
		} else {
			for i := range fields {
				validationErrors[fields[i].Field()] = fields[i].Translate(trans)
			}
		}
	}
	return validationErrors
}

// CheckJsonRequest Response API
func CheckJsonRequest(language string, err error) interface{} {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	lang := "id"
	acceptLang := language
	if i18n.IsExist(acceptLang) {
		lang = acceptLang
	}

	switch {
	case errors.As(err, &syntaxError):
		return map[string]interface{}{
			"json_body": i18n.Tr(lang, "message.jsonSyntaxError", syntaxError.Offset),
		}
	case errors.Is(err, io.ErrUnexpectedEOF):
		return map[string]interface{}{
			"json_body": i18n.Tr(lang, "message.jsonUnexpectedEof"),
		}
	case errors.As(err, &unmarshalTypeError):
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			return map[string]interface{}{
				"json_body": i18n.Tr(lang, "message.unmarshalTypeError", ute.Offset),
			}
		}
	case errors.As(err, &invalidUnmarshalError):
		if ute, ok := err.(*json.InvalidUnmarshalError); ok {
			return map[string]interface{}{
				"json_body": i18n.Tr(lang, "message.unmarshalTypeError", ute.Error()),
			}
		}
	}
	return ""
}
