package customer

import (
	"github.com/alpakih/point-of-sales/internal/domain"
	"github.com/alpakih/point-of-sales/internal/models"
	"github.com/alpakih/point-of-sales/pkg/database"
	"strings"
)

type Mapper struct {
}

func NewCustomerMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) ToCustomerResponse(customer domain.Customer) models.CustomerResponse {
	return models.CustomerResponse{
		ID:          customer.ID,
		Name:        customer.Name,
		Email:       customer.Email,
		MobilePhone: customer.MobilePhone,
	}
}

func (m *Mapper) CustomerStoreRequestToEntity(request models.CustomerStoreRequest) domain.Customer {
	return domain.Customer{
		Name:        request.Name,
		Email:       request.Email,
		MobilePhone: request.MobilePhone,
		Password:    request.Password,
	}
}

func (m *Mapper) CustomerUpdateRequestToEntity(request models.CustomerUpdateRequest, id int) domain.Customer {
	var entity domain.Customer
	entity.ID = id
	entity.Name = request.Name
	entity.Email = request.Email
	entity.MobilePhone = request.MobilePhone
	if !strings.EqualFold(request.Password, "") {
		entity.Password = request.Password
	}

	return entity
}

func (m *Mapper) ToCustomerPaginationResponse(paginator *database.Paginator) (models.Pagination, []models.CustomerResponse) {
	list := paginator.Records.(*[]domain.Customer)
	var data = make([]models.CustomerResponse, len(*list))
	for k, v := range *list {
		data[k] = models.CustomerResponse{
			ID:          v.ID,
			Name:        v.Name,
			Email:       v.Email,
			MobilePhone: v.MobilePhone,
		}
	}
	return models.BuildPaginationInfo(paginator.MaxPage, paginator.Total, paginator.PageSize, paginator.CurrentPage,
		models.BuildPaginationLinks(paginator.Links.First, paginator.Links.Prev, paginator.Links.Next, paginator.Links.Last)), data
}
