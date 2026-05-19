package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CashType string

const (
	CashIn  CashType = "in"
	CashOut CashType = "out"
)

type CashTransaction struct {
	ID          uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	COAID       uuid.UUID
	ShiftID     *uuid.UUID
	Type        CashType
	Amount      decimal.Decimal
	Description *string
	Date        time.Time
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	COA   ChartOfAccount `gorm:"foreignKey:COAID"`
	Shift *Shift         `gorm:"foreignKey:ShiftID"`
}

func (ct *CashTransaction) ToLedgerModel() Ledger {
	return Ledger{
		COAID:           ct.COAID,
		Amount:          ct.Amount,
		Notes:           ct.Description,
		ReferenceID:     ct.ID,
		ReferenceType:   LedgerCashTransaction,
		TransactionDate: ct.Date,
	}
}
