package main

import (
	customerHttpHandler "github.com/alpakih/point-of-sales/internal/customer/delivery/http"
	customerPgRepo "github.com/alpakih/point-of-sales/internal/customer/repository/pg"
	customerUCase "github.com/alpakih/point-of-sales/internal/customer/usecase"
	"github.com/alpakih/point-of-sales/internal/domain"
	"github.com/alpakih/point-of-sales/pkg/database"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"strings"
)

func main() {
	var dbSectionConfig map[string]string

	// database initialization
	if config, err := beego.AppConfig.GetSection("database"); err != nil {
		panic(err)
	} else {
		dbSectionConfig = config
	}
	db, err := database.New(database.ConfigFromEnvironment(dbSectionConfig))
	if err != nil {
		panic(err)
	}

	if err := db.Conn().AutoMigrate(&domain.Customer{}); err != nil {
		panic(err)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// init message
	languages := strings.Split(beego.AppConfig.DefaultString("lang", "en|id"), "|")
	for i := range languages {
		if err := i18n.SetMessage(languages[i], "conf/"+languages[i]+".ini"); err != nil {
			panic("Failed to set message file for l10n")
		}
	}

	customerRepository := customerPgRepo.NewCustomerPgRepository(db.Conn())
	customerUseCase := customerUCase.NewCustomerUseCase(customerRepository)
	customerHttpHandler.NewCustomerHandler(customerUseCase)

	beego.Run()
}
