package customer

import (
	"context"
)

type UseCase interface {
	StoreCustomer(ctx context.Context, request StoreRequest) (*Response, error)
	UpdateCustomer(ctx context.Context, entity UpdateRequest, id int) error
	GetCustomerByID(ctx context.Context, id int) (*Response, error)
	DeleteCustomer(ctx context.Context, id int) error
	GetCustomers(ctx context.Context, page, size int, search, order string) (*PaginationResponse, error)
}
