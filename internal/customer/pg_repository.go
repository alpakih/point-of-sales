package customer

import (
	"context"
	"github.com/alpakih/point-of-sales/internal/domain"
	"github.com/alpakih/point-of-sales/pkg/database"
)

type PgRepository interface {
	Create(ctx context.Context, entity *domain.Customer) error
	Update(ctx context.Context, entity domain.Customer) error
	FindOneCustomerByID(ctx context.Context, id int) (domain.Customer, error)
	FindCustomers(ctx context.Context, page, size int, search, order string) (*database.Paginator, error)
	CheckDuplicate(ctx context.Context, args ...interface{}) (int64, error)
	Delete(ctx context.Context, id int) error
}
