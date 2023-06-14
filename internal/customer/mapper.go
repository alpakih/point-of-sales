package customer

import (
	"github.com/alpakih/point-of-sales/internal/domain"
	"github.com/alpakih/point-of-sales/pkg/database"
	"github.com/alpakih/point-of-sales/pkg/utils"
	"strings"
)

type Mapper struct {
}

func NewCustomerMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) ToCustomerResponse(customer domain.Customer) Response {
	return Response{
		ID:          customer.ID,
		Name:        customer.Name,
		Email:       customer.Email,
		MobilePhone: customer.MobilePhone,
	}
}

func (m *Mapper) CustomerStoreRequestToEntity(request StoreRequest) domain.Customer {
	return domain.Customer{
		Name:        request.Name,
		Email:       request.Email,
		MobilePhone: request.MobilePhone,
		Password:    request.Password,
	}
}

func (m *Mapper) CustomerUpdateRequestToEntity(request UpdateRequest, id int) domain.Customer {
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

func (m *Mapper) ToCustomerPaginationResponse(paginator *database.Paginator) PaginationResponse {
	var paginationResponse PaginationResponse
	if list, ok := paginator.Records.(*[]domain.Customer); ok {
		var data = make([]Response, len(*list))
		for k, v := range *list {
			data[k] = Response{
				ID:          v.ID,
				Name:        v.Name,
				Email:       v.Email,
				MobilePhone: v.MobilePhone,
			}
		}
		paginationResponse = PaginationResponse{
			Pagination: utils.BuildPaginationInfo(
				paginator.MaxPage,
				paginator.Total,
				paginator.PageSize,
				paginator.CurrentPage,
				utils.BuildPaginationLinks(paginator.Links.First, paginator.Links.Prev, paginator.Links.Next, paginator.Links.Last)),
			Data: data,
		}
	}

	return paginationResponse
}
