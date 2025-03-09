package models

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Accounts struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AccountID string    `gorm:"uniqueIndex:idx_unique_account_id;not null" json:"account_id"`
	Balance   float64   `gorm:"not null;default:0" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Migrate(lgr *zap.SugaredLogger, db *gorm.DB) {
	err := db.AutoMigrate(&Accounts{})
	if err != nil {
		lgr.Fatalf("Migration failed: %v", err)
	}
}
