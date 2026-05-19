package dto

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LedgerResponse struct {
	ID              uuid.UUID            `json:"id"`
	COAName         string               `json:"coaName"`
	COAType         domain.COAType       `json:"coaType"`
	Amount          decimal.Decimal      `json:"amount"`
	ReferenceID     uuid.UUID            `json:"referenceID"`
	ReferenceType   domain.ReferenceType `json:"referenceType"`
	TransactionDate time.Time            `json:"transactionDate"`
	CreatedAt       time.Time            `json:"createdAt"`
}

type TransactionResponse struct {
	TransactionDate time.Time            `json:"transactionDate"`
	ReferenceType   domain.ReferenceType `json:"referenceType"`
	COAName         string               `json:"coaName"`
	COAType         domain.COAType       `json:"coaType"`
	RunningBalance  decimal.Decimal      `json:"runningBalance"`
}

type TransactionSummaryResponse struct {
	OpeningBalance decimal.Decimal       `json:"openingBalance"`
	TotalIncome    decimal.Decimal       `json:"totalIncome"`
	TotalExpense   decimal.Decimal       `json:"totalExpense"`
	Total          decimal.Decimal       `json:"total"`
	Transactions   []TransactionResponse `json:"transactions"`
	Meta           filter.Meta           `json:"meta"`
}

type CashFlowReportResponse struct {
	OpeningBalance decimal.Decimal  `json:"openingBalance"`
	TotalIncome    decimal.Decimal  `json:"totalIncome"`
	TotalExpense   decimal.Decimal  `json:"totalExpense"`
	CashFlowAmount decimal.Decimal  `json:"cashFlowAmount"`
	EndingBalance  decimal.Decimal  `json:"endingBalance"`
	Incomes        []LedgerResponse `json:"incomes"`
	Expenses       []LedgerResponse `json:"expenses"`
}
