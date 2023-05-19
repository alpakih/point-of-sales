package usecase

import (
	"context"
	"github.com/alpakih/point-of-sales/internal/customer"
	"github.com/alpakih/point-of-sales/internal/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type customerUseCase struct {
	pgRepository customer.PgRepository
}

func NewCustomerUseCase(pgRepository customer.PgRepository) customer.UseCase {
	return &customerUseCase{
		pgRepository: pgRepository,
	}
}

func (c customerUseCase) StoreCustomer(ctx context.Context, request models.CustomerStoreRequest) (*models.CustomerResponse, error) {
	var entity = customer.NewCustomerMapper().CustomerStoreRequestToEntity(request)

	// encrypt password
	if password, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		entity.Password = string(password)
	}

	// check duplicate email
	if countEmail, err := c.pgRepository.CheckDuplicate(ctx, "email =?", entity.Email); err != nil {
		return nil, err
	} else {
		if countEmail > 0 {
			return nil, models.ErrEmailAlreadyExist
		}
	}

	// check duplicate mobile phone
	countMobilePhone, err := c.pgRepository.CheckDuplicate(ctx, "mobile_phone =?", entity.MobilePhone)

	if err != nil {
		return nil, err
	}

	if countMobilePhone > 0 {
		return nil, models.ErrMobilePhoneAlreadyExist
	}

	if err := c.pgRepository.Create(ctx, &entity); err != nil {
		return nil, err
	}

	result := customer.NewCustomerMapper().ToCustomerResponse(entity)

	return &result, nil
}

func (c customerUseCase) UpdateCustomer(ctx context.Context, request models.CustomerUpdateRequest, id int) error {

	data, err := c.pgRepository.FindOneCustomerByID(ctx, id)

	if err != nil {
		return err
	}

	var entity = customer.NewCustomerMapper().CustomerUpdateRequestToEntity(request, data.ID)

	if !strings.EqualFold(entity.Password, "") {
		password, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		entity.Password = string(password)
	}

	// check duplicate email
	if countEmail, err := c.pgRepository.CheckDuplicate(ctx, "email =? and id <> ?", entity.Email, id); err != nil {
		return err
	} else {
		if countEmail > 0 {
			return models.ErrEmailAlreadyExist
		}
	}

	// check duplicate mobile phone
	if countMobilePhone, err := c.pgRepository.CheckDuplicate(ctx, "mobile_phone =? and id <> ?", entity.MobilePhone, id); err != nil {
		return err
	} else {
		if countMobilePhone > 0 {
			return models.ErrMobilePhoneAlreadyExist
		}
	}

	return c.pgRepository.Update(ctx, entity)
}

func (c customerUseCase) GetCustomerByID(ctx context.Context, id int) (*models.CustomerResponse, error) {
	data, err := c.pgRepository.FindOneCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	result := customer.NewCustomerMapper().ToCustomerResponse(data)
	return &result, nil
}

func (c customerUseCase) GetCustomers(ctx context.Context, page, size int, search, order string) (*models.Pagination, []models.CustomerResponse, error) {

	paginator, err := c.pgRepository.FindCustomers(ctx, page, size, search, order)

	if err != nil {
		return nil, nil, err
	}

	pagination, data := customer.NewCustomerMapper().ToCustomerPaginationResponse(paginator)
	return &pagination, data, nil
}

func (c customerUseCase) DeleteCustomer(ctx context.Context, id int) error {
	data, err := c.pgRepository.FindOneCustomerByID(ctx, id)

	if err != nil {
		return err
	}
	return c.pgRepository.Delete(ctx, data.ID)
}
