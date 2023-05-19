package models

type CustomerStoreRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	MobilePhone string `json:"mobile_phone" validate:"required,min=9,max=14"`
	Password    string `json:"password" validate:"required"`
}

type CustomerUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	MobilePhone string `json:"mobile_phone" validate:"required,min=9,max=14"`
	Password    string `json:"password"`
}

type CustomerResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobilePhone"`
}

type CustomerPaginationResponse struct {
	MaxPage     int64              `json:"max_page"`
	Total       int64              `json:"total"`
	PageSize    int                `json:"page_size"`
	CurrentPage int                `json:"current_page"`
	Customers   []CustomerResponse `json:"customers"`
	Links       struct {
		First string `json:"first"`
		Prev  string `json:"prev"`
		Next  string `json:"next"`
		Last  string `json:"last"`
	} `json:"links"`
}
