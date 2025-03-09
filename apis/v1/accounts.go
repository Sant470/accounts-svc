package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sant470/accounts-svc/common"
	"github.com/sant470/accounts-svc/common/errors"
	"github.com/sant470/accounts-svc/common/respond"
	"github.com/sant470/accounts-svc/services"
	apptypes "github.com/sant470/accounts-svc/types"
	"go.uber.org/zap"
)

type AccountsHlr struct {
	lgr *zap.SugaredLogger
	svc services.AccountsSvc
}

func NewAccountsHlr(lgr *zap.SugaredLogger, svc services.AccountsSvc) *AccountsHlr {
	return &AccountsHlr{lgr, svc}
}

func (impl *AccountsHlr) CreateAccounts(rw http.ResponseWriter, r *http.Request) *errors.AppError {
	var acc apptypes.AccountsRequest
	if err := common.Decode(r, &acc); err != nil {
		return errors.BadRequest(fmt.Sprintf("failed to parse requests: %s", err.Error()))
	}
	err := impl.svc.CreateAccountsSvc(&acc)
	if err != nil {
		return err
	}
	return respond.OK(rw, nil)
}

func (impl *AccountsHlr) GetAccountDetails(rw http.ResponseWriter, r *http.Request) *errors.AppError {
	accountID := chi.URLParam(r, "accountID")
	if accountID == "" {
		return errors.BadRequest("account_id is missing")
	}
	details, err := impl.svc.GetAccountDetailsSvc(accountID)
	if err != nil {
		return err
	}
	return respond.OK(rw, details)
}

func (impl *AccountsHlr) Transactions(rw http.ResponseWriter, r *http.Request) *errors.AppError {
	var req apptypes.TransactionRequest
	if err := common.Decode(r, &req); err != nil {
		return errors.BadRequest(fmt.Sprintf("failed to parse request body %s", err.Error()))
	}
	err := impl.svc.TransactionsSvc(&req)
	if err != nil {
		return err
	}
	return respond.OK(rw, "successfully transfered")
}
