package usecase

import (
	"context"
	"github.com/alpakih/point-of-sales/internal/constant"
	"github.com/alpakih/point-of-sales/internal/customer"
	"github.com/alpakih/point-of-sales/internal/customer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCustomerUseCase_StoreCustomer(t *testing.T) {
	mockCustomerRepository := new(mocks.PgRepository)
	mockDataCustomerRequest := customer.StoreRequest{
		Name:        "name",
		Email:       "email@test.com",
		MobilePhone: "087666777876",
		Password:    "123321",
	}

	t.Run("success", func(t *testing.T) {
		tempMockCustomer := mockDataCustomerRequest

		mockCustomerRepository.On("CheckDuplicate", mock.Anything, "email =?", mock.Anything).Return(int64(0), nil).Once()

		mockCustomerRepository.On("CheckDuplicate", mock.Anything, "mobile_phone =?", mock.Anything).Return(int64(0), nil).Once()

		mockCustomerRepository.On("Create", mock.Anything, mock.AnythingOfType("*domain.Customer")).Return(nil).Once()

		u := NewCustomerUseCase(mockCustomerRepository)

		data, err := u.StoreCustomer(context.TODO(), tempMockCustomer)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		mockCustomerRepository.AssertExpectations(t)
	})

	t.Run("existing-mobile-phone", func(t *testing.T) {
		tempMockCustomer := mockDataCustomerRequest

		mockCustomerRepository.On("CheckDuplicate", mock.Anything, "email =?", mock.Anything).Return(int64(0), nil).Once()

		mockCustomerRepository.On("CheckDuplicate", mock.Anything, "mobile_phone =?", mock.Anything).Return(int64(1), nil).Once()

		u := NewCustomerUseCase(mockCustomerRepository)

		data, err := u.StoreCustomer(context.TODO(), tempMockCustomer)

		assert.Error(t, err)
		assert.Nil(t, data)

		mockCustomerRepository.AssertExpectations(t)
	})

	t.Run("existing-email", func(t *testing.T) {
		tempMockCustomer := mockDataCustomerRequest

		mockCustomerRepository.On("CheckDuplicate", mock.Anything, "email =?", tempMockCustomer.Email).Return(int64(1), constant.ErrEmailAlreadyExist).Once()

		u := NewCustomerUseCase(mockCustomerRepository)

		data, err := u.StoreCustomer(context.TODO(), tempMockCustomer)

		assert.Error(t, err)
		assert.Nil(t, data)

		mockCustomerRepository.AssertExpectations(t)
	})
}
