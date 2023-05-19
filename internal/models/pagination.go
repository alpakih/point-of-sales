package models

type Pagination struct {
	MaxPage     int64           `json:"max_page"`
	Total       int64           `json:"total"`
	PageSize    int             `json:"page_size"`
	CurrentPage int             `json:"current_page"`
	Links       PaginationLinks `json:"links"`
}

type PaginationLinks struct {
	First string `json:"first"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Last  string `json:"last"`
}

func BuildPaginationLinks(first, prev, next, last string) PaginationLinks {
	return PaginationLinks{
		First: first,
		Prev:  prev,
		Next:  next,
		Last:  last,
	}
}

func BuildPaginationInfo(maxPage, total int64, pageSize, currentPage int, links PaginationLinks) Pagination {
	return Pagination{
		MaxPage:     maxPage,
		Total:       total,
		PageSize:    pageSize,
		CurrentPage: currentPage,
		Links:       links,
	}
}
