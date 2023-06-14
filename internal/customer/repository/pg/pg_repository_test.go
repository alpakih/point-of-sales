package pg

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alpakih/point-of-sales/internal/domain"
	"github.com/alpakih/point-of-sales/pkg/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"
)

func TestCustomerPgRepository_Create(t *testing.T) {

	gormDb, dbMock := utils.GetDatabaseMock("postgres")

	data := domain.Customer{
		Name:        "name",
		Email:       "email",
		MobilePhone: "mobile phone",
		Password:    "password",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO "customers" ("name","email","mobile_phone","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	queryRegex := regexp.QuoteMeta(query)

	dbMock.ExpectBegin()
	dbMock.ExpectQuery(queryRegex).WithArgs(data.Name, data.Email, data.MobilePhone, data.Password, data.CreatedAt, data.UpdatedAt).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	dbMock.ExpectCommit()

	pgRepository := NewCustomerPgRepository(gormDb)

	err := pgRepository.Create(context.TODO(), &data)
	assert.NoError(t, err)
}

func TestCustomerPgRepository_Update(t *testing.T) {
	gormDb, dbMock := utils.GetDatabaseMock("postgres")

	data := domain.Customer{
		ID:          1,
		Name:        "name",
		Email:       "email",
		MobilePhone: "mobile phone",
		Password:    "password",
		UpdatedAt:   time.Now(),
	}

	query := `UPDATE "customers" SET "name"=$1,"email"=$2,"mobile_phone"=$3,"password"=$4,"updated_at"=$5 WHERE "id" = $6`
	queryRegex := regexp.QuoteMeta(query)

	dbMock.ExpectBegin()
	dbMock.ExpectExec(queryRegex).WithArgs(data.Name, data.Email, data.MobilePhone, data.Password, utils.AnyTime{}, data.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	pgRepository := NewCustomerPgRepository(gormDb)

	err := pgRepository.Update(context.TODO(), data)
	assert.NoError(t, err)
}

func TestCustomerPgRepository_FindOneCustomerByID(t *testing.T) {
	gormDb, dbMock := utils.GetDatabaseMock("postgres")

	query := `SELECT * FROM "customers" WHERE id =$1 ORDER BY "customers"."id" LIMIT 1`
	queryRegex := regexp.QuoteMeta(query)

	dbMock.ExpectQuery(queryRegex).WithArgs(1).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "name", "email", "mobile_phone", "password", "created_at", "updated_at"}).
			AddRow(1, "name", "email", "mobile phone", "password", time.Now(), time.Now()),
	)

	pgRepository := NewCustomerPgRepository(gormDb)

	data, err := pgRepository.FindOneCustomerByID(context.TODO(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestCustomerPgRepository_FindCustomers(t *testing.T) {
	gormDb, dbMock := utils.GetDatabaseMock("postgres")

	queryCount := `SELECT count(*) FROM "customers"`
	queryRegexCount := regexp.QuoteMeta(queryCount)
	dbMock.ExpectQuery(queryRegexCount).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	query := `SELECT * FROM "customers" LIMIT 10`
	queryRegex := regexp.QuoteMeta(query)
	dbMock.ExpectQuery(queryRegex).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "email", "mobile_phone", "password", "created_at", "updated_at"}).
			AddRow(1, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(2, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(3, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(4, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(5, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(6, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(7, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(8, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(9, "name", "email", "mobile phone", "password", time.Now(), time.Now()).
			AddRow(10, "name", "email", "mobile phone", "password", time.Now(), time.Now()),
		)

	pgRepository := NewCustomerPgRepository(gormDb)

	data, err := pgRepository.FindCustomers(context.WithValue(context.TODO(), "requestCtx", &http.Request{URL: &url.URL{
		Scheme: "http",
		Host:   "localhost:8083",
		Path:   "/api/v1/customers",
	}}), 1, 10, "", "")

	assert.NoError(t, err)
	assert.NotNil(t, data.Records)

}

func TestCustomerPgRepository_Delete(t *testing.T) {
	gormDb, dbMock := utils.GetDatabaseMock("postgres")

	query := `DELETE FROM "customers" WHERE "customers"."id" = $1`
	queryRegex := regexp.QuoteMeta(query)
	dbMock.ExpectBegin()
	dbMock.ExpectExec(queryRegex).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	pgRepository := NewCustomerPgRepository(gormDb)

	err := pgRepository.Delete(context.TODO(), 1)
	assert.NoError(t, err)

}
