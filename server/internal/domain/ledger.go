package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ReferenceType string

const (
	LedgerCashTransaction ReferenceType = "cash_transactions"
	LedgerShiftSales      ReferenceType = "shifts"
)

type Ledger struct {
	ID              uuid.UUID `gorm:"primaryKey"`
	COAID           uuid.UUID
	ShiftID         *uuid.UUID
	Amount          decimal.Decimal
	Notes           *string
	ReferenceID     uuid.UUID
	ReferenceType   ReferenceType
	TransactionDate time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	COA   ChartOfAccount `gorm:"foreignKey:COAID"`
	Shift *Shift         `gorm:"foreignKey:ShiftID"`
}

type LedgerWithBalance struct {
	Ledger
	RunningBalance decimal.Decimal `gorm:"column:running_balance"`
}

type TransactionSummary struct {
	OpeningBalance decimal.Decimal `gorm:"column:opening_balance"`
	TotalIncome    decimal.Decimal `gorm:"column:total_income"`
	TotalExpense   decimal.Decimal `gorm:"column:total_expense"`
	Total          decimal.Decimal `gorm:"column:total"`
}

type CashFlowReport struct {
	OpeningBalance decimal.Decimal `gorm:"column:opening_balance"`
	TotalIncome    decimal.Decimal `gorm:"column:total_income"`
	TotalExpense   decimal.Decimal `gorm:"column:total_expense"`
	CashFlowAmount decimal.Decimal
	Total          decimal.Decimal
}

func (cfr *CashFlowReport) Calculate() {
	cfr.CashFlowAmount = cfr.TotalIncome.Sub(cfr.TotalExpense)
	cfr.Total = cfr.OpeningBalance.Add(cfr.CashFlowAmount)
}
