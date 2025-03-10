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

func TestTransactionScenarios(t *testing.T) {
	db := setupTestDB()
	store := &AccountsStoreImpl{db: db}

	db.Create(&models.Accounts{AccountID: "A1", Balance: 100})
	db.Create(&models.Accounts{AccountID: "A2", Balance: 50})

	tests := []struct {
		name       string
		fromID     string
		toID       string
		amount     float64
		wantErr    bool
		errMessage string
		finalSrc   float64
		finalDest  float64
	}{
		{
			name:      "Successful Transaction",
			fromID:    "A1",
			toID:      "A2",
			amount:    50,
			wantErr:   false,
			finalSrc:  50,
			finalDest: 100,
		},
		{
			name:       "Insufficient Funds",
			fromID:     "A1",
			toID:       "A2",
			amount:     200,
			wantErr:    true,
			errMessage: "insufficient funds or account not found",
		},
		{
			name:       "Destination Account Not Found",
			fromID:     "A1",
			toID:       "A3",
			amount:     50,
			wantErr:    true,
			errMessage: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Transaction(tt.fromID, tt.toID, tt.amount)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)

				var src, dest models.Accounts
				db.First(&src, "account_id = ?", tt.fromID)
				db.First(&dest, "account_id = ?", tt.toID)

				assert.Equal(t, tt.finalSrc, src.Balance)
				assert.Equal(t, tt.finalDest, dest.Balance)
			}
		})
	}
}

func TestConcurrentTransactions(t *testing.T) {
	db := setupTestDB()
	store := &AccountsStoreImpl{db: db}

	db.Create(&models.Accounts{AccountID: "A1", Balance: 100})
	db.Create(&models.Accounts{AccountID: "A2", Balance: 50})

	ch := make(chan error, 2)
	go func() { ch <- store.Transaction("A1", "A2", 70) }()
	go func() { ch <- store.Transaction("A1", "A2", 50) }()

	err1 := <-ch
	err2 := <-ch
	assert.NotEqual(t, err1 == nil, err2 == nil)
}
