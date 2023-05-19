package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpgx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	PostgresDriver  = "postgres"
	MysqlDriver     = "mysql"
	SqlServerDriver = "mssql"

	DefaultMaxOpenConnection     = 25
	DefaultMaxIdleConnection     = 25
	DefaultMaxLifeTimeConnection = 300
	DefaultMaxIdleTimeConnection = 300
)

var templateDsn = map[string]string{
	PostgresDriver:  "host=%s port=%s user=%s dbname=%s password=%s %s",
	MysqlDriver:     "%s:%s@(%s:%s)/%s?%s",
	SqlServerDriver: "sqlserver://%s:%s@%s:%s?database=%s&%s",
}

type Config struct {
	Driver                string
	Host                  string
	Port                  string
	Name                  string
	Username              string
	Password              string
	Options               string
	TemplateDsn           string
	Debug                 bool
	MaxOpenConnection     int
	MaxIdleConnection     int
	MaxLifeTimeConnection int
	MaxIdleTimeConnection int
	NewrelicIntegration   bool
}

func defaultDatabaseConfig() Config {

	config := Config{
		Debug:                 true,
		MaxOpenConnection:     DefaultMaxOpenConnection,
		MaxIdleConnection:     DefaultMaxIdleConnection,
		MaxLifeTimeConnection: DefaultMaxLifeTimeConnection,
		MaxIdleTimeConnection: DefaultMaxIdleTimeConnection,
		NewrelicIntegration:   false,
	}

	return config
}

func (r *Config) connectDatabase() (*gorm.DB, error) {
	var logLevel = logger.Info

	if !r.Debug {
		logLevel = logger.Silent
	}

	// default
	gormDialect, err := r.getDialect()
	if err != nil {
		return nil, err
	}

	if r.NewrelicIntegration {

		if sqlDb, err := r.newrelicDatabaseConnection(); err != nil {
			return nil, err
		} else {
			dialect, err := r.getDialectWithExistingConnection(sqlDb)
			if err != nil {
				return nil, err
			}
			gormDialect = dialect
		}
	}

	if gormDB, err := gorm.Open(
		gormDialect,
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logLevel),
		},
	); err != nil {
		return nil, err
	} else {
		dbConn, err := gormDB.DB()
		if err != nil {
			return nil, err
		}
		if err := dbConn.Ping();err!= nil{
			return nil, err
		}
		dbConn.SetMaxOpenConns(r.MaxOpenConnection)
		dbConn.SetMaxIdleConns(r.MaxIdleConnection)
		dbConn.SetConnMaxLifetime(time.Duration(r.MaxLifeTimeConnection) * time.Second)
		dbConn.SetConnMaxIdleTime(time.Duration(r.MaxIdleTimeConnection) * time.Second)
		return gormDB, nil
	}
}

func (r *Config) newrelicDatabaseConnection() (*sql.DB, error) {

	switch r.Driver {
	case "postgres":
		return sql.Open("nrpgx", r.buildDsnConnection())
	case "mysql":
		return sql.Open("nrmysql", r.buildDsnConnection())
	default:
		return nil, errors.New("unsupported driver for newrelic integration")
	}
}

func (r *Config) getDialect() (gorm.Dialector, error) {

	switch r.Driver {
	case "postgres":
		return postgres.Open(r.buildDsnConnection()), nil
	case "mssql":
		return sqlserver.Open(r.buildDsnConnection()), nil
	case "mysql":
		return mysql.Open(r.buildDsnConnection()), nil
	default:
		return nil, errors.New("unsupported driver database")
	}
}

func (r *Config) getDialectWithExistingConnection(sqlDb *sql.DB) (gorm.Dialector, error) {

	switch r.Driver {
	case "postgres":
		return postgres.New(postgres.Config{Conn: sqlDb}), nil
	case "mssql":
		return mysql.New(mysql.Config{Conn: sqlDb}), nil
	case "mysql":
		return sqlserver.New(sqlserver.Config{Conn: sqlDb}), nil
	default:
		return nil, errors.New("unsupported driver database")
	}
}

func (r *Config) buildDsnConnection() string {
	if r.Driver == "postgres" {
		return fmt.Sprintf(r.TemplateDsn, r.Host, r.Port, r.Username, r.Name, r.Password, r.Options)
	} else {
		return fmt.Sprintf(r.TemplateDsn, r.Username, r.Password, r.Host, r.Port, r.Name, r.Options)
	}
}
