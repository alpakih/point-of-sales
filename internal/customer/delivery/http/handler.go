package http

import (
	"context"
	"errors"
	"github.com/alpakih/point-of-sales/internal/customer"
	"github.com/alpakih/point-of-sales/internal/models"
	"github.com/alpakih/point-of-sales/pkg/beegoresp"
	"github.com/alpakih/point-of-sales/pkg/utils"
	"github.com/alpakih/point-of-sales/pkg/validator"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type CustomerHandler struct {
	beego.Controller
	i18n.Locale
	beegoresp.ApiResponse
	CustomerUseCase customer.UseCase
}

func NewCustomerHandler(useCase customer.UseCase) {
	handler := &CustomerHandler{
		CustomerUseCase: useCase,
	}
	beego.Router("/api/v1/customer", handler, "post:StoreCustomer")
	beego.Router("/api/v1/customer/:id", handler, "get:GetCustomerByID")
	beego.Router("/api/v1/customer/:id", handler, "put:UpdateCustomer")
	beego.Router("/api/v1/customer/:id", handler, "delete:DeleteCustomer")
	beego.Router("/api/v1/customers", handler, "get:GetCustomers")
}

func (h *CustomerHandler) Prepare() {
	h.Lang = utils.GetLangVersion(h.Ctx)
}

func (h *CustomerHandler) StoreCustomer() {
	var request models.CustomerStoreRequest

	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidJsonErrorCode, beegoresp.CheckJsonRequest(h.Lang, err))
		return
	}

	if err := validator.Validate.ValidateStruct(request); err != nil {
		h.ResponseValidationError(h.Ctx, http.StatusUnprocessableEntity, models.DataValidationErrorCode, i18n.Tr(h.Lang, "message.errorDataValidation"), err)
		return
	}

	if response, err := h.CustomerUseCase.StoreCustomer(h.Ctx.Request.Context(), request); err != nil {
		if errors.Is(err, models.ErrEmailAlreadyExist) {
			h.ResponseError(h.Ctx, http.StatusUnprocessableEntity, models.DataValidationErrorCode, i18n.Tr(h.Lang, "message.errorDataValidation"), beegoresp.DetailErrors{
				Target:      "email",
				Reason:      "duplicate",
				Description: i18n.Tr(h.Lang, "message.errorEmailAlreadyExist", request.Email),
			})
			return
		}
		if errors.Is(err, models.ErrMobilePhoneAlreadyExist) {
			h.ResponseError(h.Ctx, http.StatusUnprocessableEntity, models.DataValidationErrorCode, i18n.Tr(h.Lang, "message.errorDataValidation"), beegoresp.DetailErrors{
				Target:      "mobile_phone",
				Reason:      "duplicate",
				Description: i18n.Tr(h.Lang, "message.errorMobilePhoneAlreadyExist", request.MobilePhone),
			})
			return
		}

		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	} else {
		h.Ok(h.Ctx, response)
		return
	}
}

func (h *CustomerHandler) UpdateCustomer() {
	var request models.CustomerUpdateRequest

	id, err := strconv.Atoi(h.Ctx.Input.Param(":id"))
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorInvalidUrlParam"))
			return
		}
		if errors.Is(err, strconv.ErrRange) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorUrlParamOutOfRange"))
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	}

	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidJsonErrorCode, beegoresp.CheckJsonRequest(h.Lang, err))
		return
	}

	if err := validator.Validate.ValidateStruct(request); err != nil {
		h.ResponseValidationError(h.Ctx, http.StatusUnprocessableEntity, models.DataValidationErrorCode, i18n.Tr(h.Lang, "message.errorDataValidation"), err)
		return
	}

	if err := h.CustomerUseCase.UpdateCustomer(h.Ctx.Request.Context(), request, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.ResponseError(h.Ctx, http.StatusNotFound, models.DataNotFoundErrorCode, i18n.Tr(h.Lang, "message.errorDataNotFound"))
			return
		}

		if errors.Is(err, models.ErrEmailAlreadyExist) {
			h.ResponseError(h.Ctx, http.StatusUnprocessableEntity, models.DataValidationErrorCode, i18n.Tr(h.Lang, "message.errorDataValidation"), beegoresp.DetailErrors{
				Target:      "email",
				Reason:      "duplicate",
				Description: i18n.Tr(h.Lang, "message.errorEmailAlreadyExist", request.Email),
			})
			return
		}
		if errors.Is(err, models.ErrMobilePhoneAlreadyExist) {
			h.ResponseError(h.Ctx, http.StatusUnprocessableEntity, models.DataValidationErrorCode, i18n.Tr(h.Lang, "message.errorDataValidation"), beegoresp.DetailErrors{
				Target:      "mobile_phone",
				Reason:      "duplicate",
				Description: i18n.Tr(h.Lang, "message.errorMobilePhoneAlreadyExist", request.MobilePhone),
			})
			return
		}

		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	}
	h.Ok(h.Ctx, request)
	return
}

func (h *CustomerHandler) GetCustomers() {

	paginationQuery, err := utils.GetPaginationFromCtx(h.Ctx)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorInvalidQueryParam"))
			return
		}
		if errors.Is(err, strconv.ErrRange) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorQueryParamOutOfRange"))
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	}

	if pagination, data, err := h.CustomerUseCase.GetCustomers(
		context.WithValue(h.Ctx.Request.Context(), "requestCtx", h.Ctx.Request),
		paginationQuery.GetPage(),
		paginationQuery.GetSize(),
		paginationQuery.GetSearch(),
		paginationQuery.GetOrderBy()); err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	} else {
		h.OkWithPagination(h.Ctx, pagination, data)
		return
	}
}

func (h *CustomerHandler) GetCustomerByID() {

	id, err := strconv.Atoi(h.Ctx.Input.Param(":id"))
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorInvalidUrlParam"))
			return
		}
		if errors.Is(err, strconv.ErrRange) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorUrlParamOutOfRange"))
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	}

	if response, err := h.CustomerUseCase.GetCustomerByID(h.Ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.ResponseError(h.Ctx, http.StatusNotFound, models.DataNotFoundErrorCode, i18n.Tr(h.Lang, "message.errorDataNotFound"))
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	} else {
		h.Ok(h.Ctx, response)
		return
	}
}

func (h *CustomerHandler) DeleteCustomer() {

	id, err := strconv.Atoi(h.Ctx.Input.Param(":id"))
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorInvalidUrlParam"))
			return
		}
		if errors.Is(err, strconv.ErrRange) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, models.InvalidPathParamErrorCode, i18n.Tr(h.Lang, "message.errorUrlParamOutOfRange"))
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	}

	if err := h.CustomerUseCase.DeleteCustomer(h.Ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.ResponseError(h.Ctx, http.StatusNotFound, models.DataNotFoundErrorCode, i18n.Tr(h.Lang, "message.errorDataNotFound"))
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, models.ServerErrorCode, i18n.Tr(h.Lang, "message.errorServer"))
		return
	}
	h.Ok(h.Ctx, nil)
	return
}
