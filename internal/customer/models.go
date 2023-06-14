package customer

import "github.com/alpakih/point-of-sales/pkg/utils"

type StoreRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	MobilePhone string `json:"mobile_phone" validate:"required,min=9,max=14"`
	Password    string `json:"password" validate:"required"`
}

type UpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	MobilePhone string `json:"mobile_phone" validate:"required,min=9,max=14"`
	Password    string `json:"password"`
}

type Response struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobilePhone"`
}

type PaginationResponse struct {
	Pagination utils.Pagination
	Data       []Response
}
