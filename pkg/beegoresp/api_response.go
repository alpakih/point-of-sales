package beegoresp

import (
	"github.com/alpakih/point-of-sales/pkg/validator"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
	ut "github.com/go-playground/universal-translator"
	validatorGo "github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

type (
	ApiResponseInterface interface {
		Ok(ctx *context.Context, data interface{}) error
		OkWithPagination(ctx *context.Context, pagination interface{}, data interface{}) error
		ResponseValidationError(ctx *context.Context, httpStatus int, code, message string, err error) error
		ResponseError(ctx *context.Context, httpStatus int, code, message string, detailError ...DetailErrors) error
	}
)

type ApiResponse struct {
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	RequestId  string      `json:"request_id"`
	TimeStamp  string      `json:"timestamp"`
}

type Error struct {
	Code    string         `json:"code,omitempty"` // an application error code
	Status  string         `json:"status"`         // http status code
	Message string         `json:"message"`
	Details []DetailErrors `json:"details"`
}

type DetailErrors struct {
	Target      string `json:"target"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
}

func (r ApiResponse) Ok(ctx *context.Context, data interface{}) error {
	ctx.Output.SetStatus(http.StatusOK)
	return ctx.Resp(ApiResponse{
		Data:      data,
		RequestId: ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
		TimeStamp: time.Now().Format(time.RFC3339),
	})
}

func (r ApiResponse) OkWithPagination(ctx *context.Context, pagination interface{}, data interface{}) error {
	ctx.Output.SetStatus(http.StatusOK)
	return ctx.Resp(ApiResponse{
		Data:       data,
		Pagination: pagination,
		RequestId:  ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
		TimeStamp:  time.Now().Format(time.RFC3339),
	})
}

func (r ApiResponse) ResponseValidationError(ctx *context.Context, httpStatus int, code, message string, err error) error {
	var translator ut.Translator
	var validationErrors = make([]DetailErrors, 0)

	ctx.Output.SetStatus(httpStatus)
	lang := "id"
	acceptLang := ctx.Request.Header.Get("Accept-Language")

	if acceptLang != "" && i18n.IsExist(acceptLang) {
		lang = acceptLang
	}
	if trans, found := validator.Validate.GetTranslator(lang); found {
		translator = trans
	}

	if fieldType, ok := err.(*validatorGo.InvalidValidationError); ok {
		validationErrors = append(validationErrors, DetailErrors{
			Target:      fieldType.Type.Name(),
			Reason:      fieldType.Type.Kind().String(),
			Description: fieldType.Error(),
		})
	}

	if fieldErrors, ok := err.(validatorGo.ValidationErrors); ok {
		for i := range fieldErrors {
			validationErrors = append(validationErrors, DetailErrors{
				Target:      fieldErrors[i].Field(),
				Reason:      fieldErrors[i].Tag(),
				Description: fieldErrors[i].Translate(translator),
			})
		}
	}

	return ctx.Resp(ApiResponse{
		Error: Error{
			Code:    code,
			Status:  strconv.Itoa(httpStatus),
			Message: message,
			Details: validationErrors,
		},
		RequestId: ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
		TimeStamp: time.Now().Format(time.RFC3339),
	})
}

func (r ApiResponse) ResponseError(ctx *context.Context, httpStatus int, code, message string, detailErrors ...DetailErrors) error {
	ctx.Output.SetStatus(httpStatus)

	var details = make([]DetailErrors, len(detailErrors))

	for k, v := range detailErrors {
		details[k] = DetailErrors{
			Target:      v.Target,
			Reason:      v.Reason,
			Description: v.Description,
		}
	}
	return ctx.Resp(ApiResponse{
		Error: Error{
			Code:    code,
			Status:  strconv.Itoa(httpStatus),
			Message: message,
			Details: details,
		},
		RequestId: ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
		TimeStamp: time.Now().Format(time.RFC3339),
	})
}
