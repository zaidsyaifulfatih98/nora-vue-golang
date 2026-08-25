package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type FinanceEntryType string

const (
	FinanceIncome  FinanceEntryType = "INCOME"
	FinanceExpense FinanceEntryType = "EXPENSE"
)

// FinanceEntry is defined for schema parity with the original Prisma model,
// but per product decision, finance data is stored in Google Sheets, not
// Postgres. This table remains unused by the finance module.
type FinanceEntry struct {
	BaseModel
	Type        FinanceEntryType `gorm:"type:text;not null" json:"type"`
	Category    string           `gorm:"size:100;not null" json:"category"`
	Amount      decimal.Decimal  `gorm:"type:decimal;not null" json:"amount"`
	Description string           `gorm:"type:text" json:"description"`
	Date        time.Time        `gorm:"not null" json:"date"`
	CreatedByID string           `gorm:"type:uuid;not null" json:"createdById"`
	CreatedBy   User             `gorm:"foreignKey:CreatedByID" json:"createdBy,omitempty"`
}

func (FinanceEntry) TableName() string { return "finance_entries" }
