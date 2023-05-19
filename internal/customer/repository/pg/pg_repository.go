package pg

import (
	"context"
	"fmt"
	"github.com/alpakih/point-of-sales/internal/customer"
	"github.com/alpakih/point-of-sales/internal/domain"
	"github.com/alpakih/point-of-sales/pkg/database"
	"github.com/alpakih/point-of-sales/pkg/utils"
	"gorm.io/gorm"
	"net/http"
)

type customerPgRepository struct {
	db *gorm.DB
}

func NewCustomerPgRepository(db *gorm.DB) customer.PgRepository {
	return &customerPgRepository{
		db: db,
	}
}

func (c customerPgRepository) Create(ctx context.Context, entity *domain.Customer) error {
	return c.db.WithContext(ctx).Create(entity).Error
}

func (c customerPgRepository) Update(ctx context.Context, entity domain.Customer) error {
	return c.db.WithContext(ctx).Updates(&entity).Error
}

func (c customerPgRepository) FindCustomers(ctx context.Context, page, size int, search, order string) (*database.Paginator, error) {
	var entities []domain.Customer
	db := c.db
	fields := utils.GetListValueFromTagStruct(domain.Customer{}, "qsearch")
	if search != "" {
		for i := range fields {
			db = db.Or(fmt.Sprintf("%s ILIKE ?", fields[i]), "%"+search+"%")
		}
	}
	if order != "" {
		if utils.ItemExists(fields, order) {
			db = db.Order(order)
		}
	}

	paginator := database.NewPaginator(db, ctx.Value("requestCtx").(*http.Request), page, size, &entities)

	return paginator, paginator.Find(ctx).Error
}

func (c customerPgRepository) FindOneCustomerByID(ctx context.Context, id int) (domain.Customer, error) {
	var entity domain.Customer
	err := c.db.WithContext(ctx).First(&entity, "id =?", id).Error
	return entity, err
}

func (c customerPgRepository) CheckDuplicate(ctx context.Context, args ...interface{}) (int64, error) {

	var count int64

	db := c.db.WithContext(ctx).Model(&domain.Customer{})

	if args != nil {
		db.Where(args[0], args[1:]...)
	}

	return count, db.Count(&count).Error
}

func (c customerPgRepository) Delete(ctx context.Context, id int) error {
	return c.db.WithContext(ctx).Delete(&domain.Customer{}, id).Error
}
