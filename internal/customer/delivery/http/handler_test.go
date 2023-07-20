package http

import (
	"encoding/json"
	"fmt"
	"github.com/alpakih/point-of-sales/internal/customer"
	"github.com/alpakih/point-of-sales/internal/customer/mocks"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	fmt.Println(file)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator)+".."+string(filepath.Separator)+".."+string(filepath.Separator)+".."+string(filepath.Separator))))
	fmt.Println(appPath)

	beego.TestBeegoInit(appPath)
}

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

		h := beego.NewControllerRegister()

		handler := &CustomerHandler{
			Locale:          i18n.Locale{Lang: "id"},
			CustomerUseCase: mockUCase,
		}

		h.Add("/api/v1/customer", handler, beego.WithRouterMethods(handler, "post:StoreCustomer"))

		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUCase.AssertExpectations(t)
	}
}
