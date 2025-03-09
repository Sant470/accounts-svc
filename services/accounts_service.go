package services

import (
	"fmt"

	"github.com/sant470/accounts-svc/common/errors"
	"github.com/sant470/accounts-svc/store"
	apptypes "github.com/sant470/accounts-svc/types"
	"go.uber.org/zap"
)

type AccountsSvc interface {
	CreateAccountsSvc(accounts *apptypes.AccountsRequest) *errors.AppError
	GetAccountDetailsSvc(accountID string) (interface{}, *errors.AppError)
	TransactionsSvc(req *apptypes.TransactionRequest) *errors.AppError
}

type AccountsSvcImpl struct {
	lgr   *zap.SugaredLogger
	store store.AccountsStore
}

func NewAccountsSvc(lgr *zap.SugaredLogger, accountsStore store.AccountsStore) *AccountsSvcImpl {
	return &AccountsSvcImpl{lgr, accountsStore}
}

func (impl *AccountsSvcImpl) CreateAccountsSvc(accounts *apptypes.AccountsRequest) *errors.AppError {
	err := impl.store.CreateAccount(accounts)
	if err != nil {
		msg := fmt.Sprintf("error creating account :%s", err.Error())
		impl.lgr.Debug(msg)
		return errors.InternalServerError(msg)
	}
	return nil
}

func (impl *AccountsSvcImpl) GetAccountDetailsSvc(accountID string) (interface{}, *errors.AppError) {
	details, err := impl.store.AccountDetails(accountID)
	if err != nil {
		msg := fmt.Sprintf("error getting account details:%s", err.Error())
		impl.lgr.Debug(msg)
		return nil, errors.InternalServerError(msg)
	}
	return details, nil
}

func (impl *AccountsSvcImpl) TransactionsSvc(req *apptypes.TransactionRequest) *errors.AppError {
	// add input validation logic here
	fromID, toID, amount := req.SourceAccountID, req.DestinationAccountID, req.Amount
	err := impl.store.Transaction(fromID, toID, amount)
	if err != nil {
		msg := fmt.Sprintf("error transfering money: %s", err.Error())
		impl.lgr.Debug(msg)
		return errors.InternalServerError(msg)
	}
	return nil
}
