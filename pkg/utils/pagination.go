package utils

import (
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"math"
	"strconv"
)

const (
	defaultSize = 10
)

type PaginationQuery struct {
	Size    int    `json:"size,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
	Search  string `json:"search,omitempty"`
}

func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

func (q *PaginationQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

func (q *PaginationQuery) SetSearch(searchQuery string) {
	q.Search = searchQuery
}

func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

func (q *PaginationQuery) GetOrderBy() string {
	return q.OrderBy
}

func (q *PaginationQuery) GetSearch() string {
	return q.Search
}

func (q *PaginationQuery) GetPage() int {
	return q.Page
}

func (q *PaginationQuery) GetSize() int {
	return q.Size
}

func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s&search=%s", q.GetPage(), q.GetSize(), q.GetOrderBy(), q.GetSearch())
}

func GetPaginationFromCtx(c *context.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.Input.Query("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.Input.Query("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.Input.Query("orderBy"))
	q.SetSearch(c.Input.Query("search"))

	return q, nil
}

func GetTotalPages(totalCount int, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

func GetHasMore(currentPage int, totalCount int, pageSize int) bool {
	return currentPage < totalCount/pageSize
}
