package database

import (
	"context"
	"gorm.io/gorm/clause"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

// Paginator structure containing pagination information and result records.
// Can be sent to the client directly.
type Paginator struct {
	DB      *gorm.DB      `json:"-"`
	Request *http.Request `json:"-"`

	Records interface{} `json:"records"`

	rawQuery          string
	rawQueryVars      []interface{}
	rawCountQuery     string
	rawCountQueryVars []interface{}

	MaxPage        int64 `json:"max_page"`
	Total          int64 `json:"total"`
	PageSize       int   `json:"page_size"`
	CurrentPage    int   `json:"current_page"`
	Links          Links `json:"links"`
	loadedPageInfo bool
}

type Links struct {
	First string `json:"first"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Last  string `json:"last"`
}

func paginateScope(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// NewPaginator create a new Paginator.
//
// Given DB transaction can contain clauses already, such as WHERE, if you want to
// filter results.
//
//  articles := []model.Article{}
//  tx := database.Conn().Where("title LIKE ?", "%"+sqlutil.EscapeLike(search)+"%")
//  paginator := database.NewPaginator(tx, page, pageSize, &articles)
//  result := paginator.Find()
//  if response.HandleDatabaseError(result) {
//      response.JSON(http.StatusOK, paginator)
//  }
//
func NewPaginator(db *gorm.DB, request *http.Request, page, pageSize int, dest interface{}) *Paginator {
	return &Paginator{
		DB:          db,
		Request:     request,
		CurrentPage: page,
		PageSize:    pageSize,
		Records:     dest,
	}
}

// Raw set a raw SQL query and count query.
// The Paginator will execute the raw queries instead of automatically creating them.
// The raw query should not contain the "LIMIT" and "OFFSET" clauses, they will be added automatically.
// The count query should return a single number (`COUNT(*)` for example).
func (p *Paginator) Raw(query string, vars []interface{}, countQuery string, countVars []interface{}) *Paginator {
	p.rawQuery = query
	p.rawQueryVars = vars
	p.rawCountQuery = countQuery
	p.rawCountQueryVars = vars
	return p
}

// UpdatePageInfo executes count request to calculate the `Total` and `MaxPage`.
func (p *Paginator) UpdatePageInfo(ctx context.Context) {
	count := int64(0)
	db := p.DB.WithContext(ctx).Session(&gorm.Session{})
	prevPreloads := db.Statement.Preloads
	if len(prevPreloads) > 0 {
		db.Statement.Preloads = map[string][]interface{}{}
		defer func() {
			db.Statement.Preloads = prevPreloads
		}()
	}
	var err error
	if p.rawCountQuery != "" {
		err = db.Raw(p.rawCountQuery, p.rawCountQueryVars...).Scan(&count).Error
	} else {
		err = db.Model(p.Records).Count(&count).Error
	}
	if err != nil {
		panic(err)
	}
	p.Total = count
	p.MaxPage = int64(math.Ceil(float64(count) / float64(p.PageSize)))
	if p.MaxPage == 0 {
		p.MaxPage = 1
	}
	p.Links.First = p.PageLinkFirst()
	p.Links.Next = p.PageLinkNext()
	p.Links.Prev = p.PageLinkPrev()
	p.Links.Last = p.PageLinkLast()

	p.loadedPageInfo = true
}

// Find requests page information (total records and max page) and
// executes the transaction. The Paginate struct is updated automatically, as
// well as the destination slice given in NewPaginator().
func (p *Paginator) Find(ctx context.Context) *gorm.DB {
	if !p.loadedPageInfo {
		p.UpdatePageInfo(ctx)
	}
	if p.rawQuery != "" {
		return p.rawStatement(ctx).Scan(p.Records)
	}
	return p.DB.Scopes(paginateScope(p.CurrentPage, p.PageSize)).Find(p.Records)
}

func (p *Paginator) rawStatement(ctx context.Context) *gorm.DB {
	offset := (p.CurrentPage - 1) * p.PageSize
	db := p.DB.WithContext(ctx).Raw(p.rawQuery, p.rawQueryVars...)

	db.Statement.SQL.WriteString(" ")

	if db.Dialector.Name() == "sqlserver" {
		if db.Statement.Schema != nil && db.Statement.Schema.PrioritizedPrimaryField != nil {
			db.Statement.SQL.WriteString("ORDER BY ")
			db.Statement.WriteQuoted(db.Statement.Schema.PrioritizedPrimaryField.DBName)
			db.Statement.SQL.WriteByte(' ')
		} else {
			db.Statement.SQL.WriteString("ORDER BY (SELECT NULL) ")
		}

		if p.CurrentPage > 0 {
			db.Statement.SQL.WriteString("OFFSET ")
			db.Statement.SQL.WriteString(strconv.Itoa(offset))
			db.Statement.SQL.WriteString(" ROWS")
		}

		if p.PageSize > 0 {
			if p.CurrentPage == 0 {
				db.Statement.SQL.WriteString("OFFSET 0 ROW")
			}
			db.Statement.SQL.WriteString(" FETCH NEXT ")
			db.Statement.SQL.WriteString(strconv.Itoa(p.PageSize))
			db.Statement.SQL.WriteString(" ROWS ONLY")
		}
	} else {
		clause.Limit{Limit: &p.PageSize, Offset: offset}.Build(db.Statement)
	}

	return db
}

func (p *Paginator) PageLink(page int) string {
	link, err := url.ParseRequestURI(p.Request.URL.String())
	if err != nil {
		panic(err)
	}
	values := link.Query()
	if page < 1 {
		values.Set("page", strconv.Itoa(page))
	} else {
		values.Set("page", strconv.Itoa(page))
	}

	link.RawQuery = values.Encode()
	return link.String()
}

// PageLinkPrev Returns URL to the previous page.
func (p *Paginator) PageLinkPrev() (link string) {
	if p.HasPrev() {
		link = p.PageLink(p.CurrentPage - 1)
	}
	return
}

// PageLinkNext Returns URL to the next page.
func (p *Paginator) PageLinkNext() (link string) {
	if p.HasNext() {
		link = p.PageLink(p.CurrentPage + 1)
	}
	return
}

// PageLinkFirst Returns URL to the first page.
func (p *Paginator) PageLinkFirst() (link string) {
	return p.PageLink(1)
}

// PageLinkLast Returns URL to the last page.
func (p *Paginator) PageLinkLast() (link string) {
	return p.PageLink(int(p.MaxPage))
}

// HasPrev Returns true if the current page has a predecessor.
func (p *Paginator) HasPrev() bool {
	return p.CurrentPage > 1
}

// HasNext Returns true if the current page has a successor.
func (p *Paginator) HasNext() bool {
	return p.CurrentPage < int(p.MaxPage)
}
