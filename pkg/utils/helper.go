package utils

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
	"strings"
	"time"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// GetLangVersion sets site language version.
func GetLangVersion(ctx *context.Context) string {
	// 1. Check URL arguments.
	lang := ctx.Input.Query("lang")

	// Check again in case someone modifies on purpose.
	if !i18n.IsExist(lang) {
		lang = ""
	}

	// 2. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := ctx.Request.Header.Get("Accept-Language")
		if i18n.IsExist(al) {
			lang = al
		}
	}

	// 3. Default language is Indonesia.
	if len(lang) == 0 {
		lang = "id"
	}

	// Set language properties.
	return lang
}

func GetListValueFromTagStruct(obj interface{}, tag string) []string {
	var field []string
	v := reflect.ValueOf(obj)
	for i := 0; i < v.Type().NumField(); i++ {
		if !strings.EqualFold(v.Type().Field(i).Tag.Get(tag), "-") {
			field = append(field, v.Type().Field(i).Tag.Get(tag))
		}
	}
	return field
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func GetDatabaseMock(driver string) (*gorm.DB, sqlmock.Sqlmock) {
	if sqlDB, mock, err := sqlmock.New(); err != nil {
		return nil, nil
	} else {
		var director gorm.Dialector

		switch driver {

		case "postgres":
			director = postgres.New(postgres.Config{Conn: sqlDB})
		case "mysql":
			director = mysql.New(mysql.Config{Conn: sqlDB})
		case "mssql":
			director = sqlserver.New(sqlserver.Config{Conn: sqlDB})
		default:
			director = postgres.New(postgres.Config{Conn: sqlDB})
		}

		gormDB, _ := gorm.Open(director, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		return gormDB, mock
	}
}
