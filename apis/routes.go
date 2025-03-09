package apis

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	v1 "github.com/sant470/accounts-svc/apis/v1"
	"github.com/sant470/accounts-svc/common"
)

func InitRegistrationHlr(r *chi.Mux, ah *v1.AccountsHlr) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/accounts", common.Handler(ah.CreateAccounts))
		r.Method(http.MethodGet, "/accounts/{accountID}", common.Handler(ah.GetAccountDetails))
		r.Method(http.MethodPost, "/transactions", common.Handler(ah.Transactions))
	})
}
