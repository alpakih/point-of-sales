package customer

import (
	"context"
	"github.com/alpakih/point-of-sales/internal/models"
)

type UseCase interface {
	StoreCustomer(ctx context.Context, request models.CustomerStoreRequest) (*models.CustomerResponse, error)
	UpdateCustomer(ctx context.Context, entity models.CustomerUpdateRequest, id int) error
	GetCustomerByID(ctx context.Context, id int) (*models.CustomerResponse, error)
	DeleteCustomer(ctx context.Context, id int) error
	GetCustomers(ctx context.Context, page, size int, search, order string) (*models.Pagination, []models.CustomerResponse, error)
}
