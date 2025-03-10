package store

import (
	"testing"

	"github.com/sant470/accounts-svc/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Accounts{})
	return db
}

func TestTransactionSuccess(t *testing.T) {
	db := setupTestDB()
	store := &AccountsStoreImpl{db: db}

	db.Create(&models.Accounts{AccountID: "A1", Balance: 100})
	db.Create(&models.Accounts{AccountID: "A2", Balance: 50})

	err := store.Transaction("A1", "A2", 50)
	assert.Nil(t, err)

	var src, dest models.Accounts
	db.First(&src, "account_id = ?", "A1")
	db.First(&dest, "account_id = ?", "A2")
	assert.Equal(t, float64(50), src.Balance)
	assert.Equal(t, float64(100), dest.Balance)
}

func TestInsufficientFunds(t *testing.T) {
	db := setupTestDB()
	store := &AccountsStoreImpl{db: db}

	db.Create(&models.Accounts{AccountID: "A1", Balance: 30})
	db.Create(&models.Accounts{AccountID: "A2", Balance: 50})

	err := store.Transaction("A1", "A2", 50)
	assert.NotNil(t, err)
	assert.Equal(t, "insufficient funds or account not found", err.Error())
}

func TestDestinationAccountNotFound(t *testing.T) {
	db := setupTestDB()
	store := &AccountsStoreImpl{db: db}

	db.Create(&models.Accounts{AccountID: "A1", Balance: 100})

	err := store.Transaction("A1", "A3", 50)
	assert.NotNil(t, err)
	assert.Equal(t, "account not found", err.Error())
}

func TestConcurrentTransactions(t *testing.T) {
	db := setupTestDB()
	store := &AccountsStoreImpl{db: db}

	db.Create(&models.Accounts{AccountID: "A1", Balance: 100})
	db.Create(&models.Accounts{AccountID: "A2", Balance: 50})

	ch := make(chan error)
	go func() { ch <- store.Transaction("A1", "A2", 70) }()
	go func() { ch <- store.Transaction("A1", "A2", 50) }()

	err1 := <-ch
	err2 := <-ch

	assert.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil))
}
