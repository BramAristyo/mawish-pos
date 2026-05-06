package dto

import (
	"time"

	"github.com/BramAristyo/saas-pos-core/server/internal/domain"
	"github.com/BramAristyo/saas-pos-core/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CashTransactionResponse struct {
	ID          uuid.UUID       `json:"id"`
	COAID       uuid.UUID       `json:"coaId"`
	CoaName     string          `json:"coaName"`
	Type        domain.CashType `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Description *string         `json:"description"`
	Date        string          `json:"date"`
	CreatedAt   string          `json:"createdAt"`
}

type CashTransactionResponsePagination struct {
	Data []CashTransactionResponse `json:"data"`
	Meta filter.Meta               `json:"meta"`
}

type CashTransactionRequest struct {
	CoaID       uuid.UUID       `json:"coaId" binding:"required,uuid"`
	ShiftID     *uuid.UUID      `json:"shiftId"`
	Type        domain.CashType `json:"type" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required,min=0"`
	Description *string         `json:"description"`
	Date        string          `json:"date" binding:"required"`
}

func ToCashTransactionModel(ctr *CashTransactionRequest) (domain.CashTransaction, error) {
	parsedDate, err := time.Parse("2006-01-02", ctr.Date)
	if err != nil {
		return domain.CashTransaction{}, err
	}
	return domain.CashTransaction{
		COAID:       ctr.CoaID,
		ShiftID:     ctr.ShiftID,
		Type:        ctr.Type,
		Amount:      ctr.Amount,
		Description: ctr.Description,
		Date:        parsedDate,
	}, nil
}

func ToCashTransactionResponse(ct *domain.CashTransaction) CashTransactionResponse {
	return CashTransactionResponse{
		ID:          ct.ID,
		CoaName:     ct.COA.Name,
		COAID:       ct.COA.ID,
		Type:        ct.Type,
		Amount:      ct.Amount,
		Description: ct.Description,
		Date:        ct.Date.Format("2006-01-02 15:04:05"),
		CreatedAt:   ct.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToCashTransactionPagination(
	cts []domain.CashTransaction,
	f filter.PaginationWithInputFilter,
	totalRows int64,
) CashTransactionResponsePagination {
	res := make([]CashTransactionResponse, len(cts))
	for i, ct := range cts {
		res[i] = ToCashTransactionResponse(&ct)
	}

	return CashTransactionResponsePagination{
		Data: res,
		Meta: f.ToMeta(totalRows),
	}
}
