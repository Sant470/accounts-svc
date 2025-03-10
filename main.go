package main

import (
	"net/http"

	"github.com/sant470/accounts-svc/apis"
	handlers "github.com/sant470/accounts-svc/apis/v1"
	"github.com/sant470/accounts-svc/config"
	"github.com/sant470/accounts-svc/models"
	"github.com/sant470/accounts-svc/services"
	"github.com/sant470/accounts-svc/store"
)

func main() {
	conf := config.GetAppConfig("config.yaml", "./")
	lgr := config.GetConsoleLogger()
	db := config.GetDBConn(lgr, conf.DB)
	models.Migrate(lgr, db)
	accountsStore := store.NewAccountStore(lgr, db)
	accountsSvc := services.NewAccountsSvc(lgr, accountsStore)
	accountsHlr := handlers.NewAccountsHlr(lgr, accountsSvc)
	router := config.InitRouters()
	apis.InitAccountsRoutes(router, accountsHlr)
	http.ListenAndServe("localhost:8000", router)
}
