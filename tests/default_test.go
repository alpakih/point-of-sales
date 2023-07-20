package test

import (
	"encoding/json"
	"fmt"
	"github.com/alpakih/point-of-sales/internal/customer"
	cHandler "github.com/alpakih/point-of-sales/internal/customer/delivery/http"
	"github.com/alpakih/point-of-sales/internal/customer/mocks"
	"github.com/beego/i18n"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	fmt.Println(file)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	fmt.Println(appPath)
	beego.TestBeegoInit(appPath)
}

// TestCustomerStore is a sample to run an endpoint test
func TestCustomerStore(t *testing.T) {

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

		handler := &cHandler.CustomerHandler{
			Locale:          i18n.Locale{Lang: "id"},
			CustomerUseCase: mockUCase,
		}

		h.Add("/api/v1/customer", handler, beego.WithRouterMethods(handler, "post:StoreCustomer"))

		h.ServeHTTP(w, r)

		Convey("Subject: Test Station Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("The Result Should Not Be Empty", func() {
				So(w.Body.Len(), ShouldBeGreaterThan, 0)
			})
		})
		mockUCase.AssertExpectations(t)
	}
}
