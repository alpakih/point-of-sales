package http

import (
	"encoding/json"
	"github.com/alpakih/point-of-sales/internal/customer"
	"github.com/alpakih/point-of-sales/internal/customer/mocks"
	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCustomerHandler_StoreCustomer(t *testing.T) {

	mockDataCustomer := customer.StoreRequest{
		Name:        "Test",
		Email:       "email@test.com",
		MobilePhone: "087666777656",
		Password:    "123123",
	}

	mockCustomerEntity := customer.NewCustomerMapper().CustomerStoreRequestToEntity(mockDataCustomer)

	mockCustomerResponse := customer.NewCustomerMapper().ToCustomerResponse(mockCustomerEntity)

	mockUCase := new(mocks.UseCase)
	mockUCase.On("StoreCustomer", mock.Anything, mock.AnythingOfType("customer.StoreRequest")).Return(&mockCustomerResponse, nil)

	if bodyJson, err := json.Marshal(mockDataCustomer); err != nil {
		assert.NoError(t, err)
	} else {
		r, err := http.NewRequest("POST", "/api/v1/customer", strings.NewReader(string(bodyJson)))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		handler := CustomerHandler{
			Controller: beego.Controller{
				Ctx: &beegoContext.Context{
					Request: r,
					ResponseWriter: &beegoContext.Response{
						ResponseWriter: w,
					},
				},
				Data: map[interface{}]interface{}{},
			},
			Locale: i18n.Locale{
				Lang: "id",
			},
			CustomerUseCase: mockUCase,
		}

		handler.Ctx.Input = &beegoContext.BeegoInput{
			Context:     handler.Ctx,
			RequestBody: bodyJson,
		}
		handler.Ctx.Output = &beegoContext.BeegoOutput{
			Context: handler.Ctx,
		}

		handler.StoreCustomer()

		assert.Equal(t, http.StatusOK, w.Code)
		mockUCase.AssertExpectations(t)
	}
}
