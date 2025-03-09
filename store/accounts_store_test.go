package store

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, *AccountsStoreImpl) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	lgr := zap.NewNop().Sugar()
	accountStore := NewAccountStore(lgr, gormDB)
	return mock, accountStore
}

func TestTransactionSuccess(t *testing.T) {
	mock, store := setupTestDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts").
		WithArgs("101", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE accounts").
		WithArgs("102", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := store.Transaction("101", "105", 100.0)
	fmt.Println("error: ", err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_InsufficientFunds(t *testing.T) {
	mock, store := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts").
		WithArgs("101", 1000.0).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Simulate insufficient funds
	mock.ExpectRollback()

	err := store.Transaction("101", "102", 1000.0)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds or account not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_MissingSender(t *testing.T) {
	mock, store := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts").
		WithArgs("missing_user", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Sender not found
	mock.ExpectRollback()

	err := store.Transaction("missing_user", "user2", 100.0)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds or account not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_MissingReceiver(t *testing.T) {
	mock, store := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts").
		WithArgs("101", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Deduction success
	mock.ExpectExec("UPDATE accounts").
		WithArgs("missing_user", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Receiver not found
	mock.ExpectRollback()

	err := store.Transaction("101", "missing_user", 100.0)
	assert.Error(t, err)
	assert.Equal(t, "account not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_Atomicity(t *testing.T) {
	mock, store := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts").
		WithArgs("101", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE accounts").
		WithArgs("102", 100.0).
		WillReturnError(errors.New("DB error")) // Simulate DB failure
	mock.ExpectRollback()

	err := store.Transaction("101", "102", 100.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB error") // Ensure error is propagated
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_ConcurrentTransfers(t *testing.T) {
	mock, store := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts").
		WithArgs("101", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE accounts").
		WithArgs("102", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts").
		WithArgs("101", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE accounts").
		WithArgs("102", 100.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	concurrentTransfers := func() {
		err := store.Transaction("101", "102", 100.0)
		assert.NoError(t, err)
	}

	// Run transactions in parallel
	go concurrentTransfers()
	go concurrentTransfers()

	// Allow goroutines to execute
	time.Sleep(200 * time.Millisecond)

	assert.NoError(t, mock.ExpectationsWereMet())
}
