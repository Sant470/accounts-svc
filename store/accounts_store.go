package store

import (
	"context"
	"errors"
	"time"

	"github.com/sant470/accounts-svc/models"
	apptypes "github.com/sant470/accounts-svc/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AccountsStore interface {
	CreateAccount(accounts *apptypes.AccountsRequest) error
	AccountDetails(accountID string) (interface{}, error)
	Transaction(fromID, toID string, amount float64) error
}

type AccountsStoreImpl struct {
	lgr *zap.SugaredLogger
	db  *gorm.DB
}

func NewAccountStore(lgr *zap.SugaredLogger, db *gorm.DB) *AccountsStoreImpl {
	return &AccountsStoreImpl{lgr, db}
}

func deductBalance(tx *gorm.DB, accountID string, amount float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := tx.WithContext(ctx).Model(&models.Accounts{}).
		Where("account_id = ? AND balance >= ?", accountID, amount).
		Update("balance", gorm.Expr("balance - ?", amount))

	if res.RowsAffected == 0 {
		return errors.New("insufficient funds or account not found")
	}
	return nil
}

func addBalance(tx *gorm.DB, accountID string, amount float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := tx.WithContext(ctx).Model(&models.Accounts{}).
		Where("account_id = ?", accountID).
		Update("balance", gorm.Expr("balance + ?", amount))
	if res.RowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (impl *AccountsStoreImpl) CreateAccount(accounts *apptypes.AccountsRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	acc := models.Accounts{AccountID: accounts.AccountID, Balance: accounts.InitialBalance}
	res := impl.db.WithContext(ctx).Create(&acc)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (impl *AccountsStoreImpl) AccountDetails(accountID string) (interface{}, error) {
	var account models.Accounts
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := impl.db.WithContext(ctx).Where("account_id = ?", accountID).First(&account).Error
	return account, err
}

func (impl *AccountsStoreImpl) Transaction(fromID, toID string, amount float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return impl.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := deductBalance(tx, fromID, amount); err != nil {
			return err
		}
		return addBalance(tx, toID, amount)
	})
}
