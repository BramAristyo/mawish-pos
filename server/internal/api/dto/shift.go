package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ShiftExpenseResponse struct {
	ID          uuid.UUID              `json:"id"`
	ShiftID     uuid.UUID              `json:"shiftId"`
	COAID       uuid.UUID              `json:"coaId"`
	COA         ChartOfAccountResponse `json:"coa"`
	Amount      decimal.Decimal        `json:"amount"`
	Description *string                `json:"description"`
	CreatedAt   string                 `json:"createdAt"`
}

type ShiftResponse struct {
	ID            uuid.UUID              `json:"id"`
	OpenedBy      uuid.UUID              `json:"openedBy"`
	OpenedByUser  UserResponse           `json:"openedByUser"`
	ClosedBy      *uuid.UUID             `json:"closedBy,omitempty"`
	ClosedByUser  *UserResponse          `json:"closedByUser,omitempty"`
	OpeningCash   decimal.Decimal        `json:"openingCash"`
	ClosingCash   *decimal.Decimal       `json:"closingCash,omitempty"`
	Notes         *string                `json:"notes"`
	OpenedAt      string                 `json:"openedAt"`
	ClosedAt      *string                `json:"closedAt,omitempty"`
	CreatedAt     string                 `json:"createdAt"`
	ShiftExpenses []ShiftExpenseResponse `json:"shiftExpenses,omitempty"`
}

type ShiftResponsePagination struct {
	Data []ShiftResponse `json:"data"`
	Meta filter.Meta     `json:"meta"`
}

type OpenShiftRequest struct {
	OpeningCash decimal.Decimal `json:"openingCash" binding:"required,gt=0"`
	Notes       *string         `json:"notes" binding:"omitempty,max=255"`
}

type CloseShiftRequest struct {
	ClosingCash decimal.Decimal `json:"closingCash" binding:"required,gt=0"`
	Notes       *string         `json:"notes" binding:"omitempty,max=255"`
}

type ShiftExpenseRequest struct {
	ID          *string         `json:"id" binding:"omitempty,uuid"`
	COAID       uuid.UUID       `json:"coaId" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required,gt=0"`
	Description *string         `json:"description" binding:"omitempty,max=255"`
}

type UpsertShiftExpensesRequest struct {
	Expenses []ShiftExpenseRequest `json:"expenses" binding:"required,min=1,dive"`
}

func ToShiftResponse(s *domain.Shift) ShiftResponse {
	var closedAt *string
	if s.ClosedAt != nil {
		str := s.ClosedAt.Format("2006-01-02 15:04:05")
		closedAt = &str
	}

	var closedByUser *UserResponse
	if s.ClosedByUser != nil {
		res := ToUserResponse(s.ClosedByUser)
		closedByUser = &res
	}

	return ShiftResponse{
		ID:           s.ID,
		OpenedBy:     s.OpenedBy,
		OpenedByUser: ToUserResponse(&s.OpenedByUser),
		ClosedBy:     s.ClosedBy,
		ClosedByUser: closedByUser,
		OpeningCash:  s.OpeningCash,
		ClosingCash:  s.ClosingCash,
		Notes:        s.Notes,
		OpenedAt:     s.OpenedAt.Format("2006-01-02 15:04:05"),
		ClosedAt:     closedAt,
		CreatedAt:    s.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToShiftResponsePagination(s []ShiftResponse, p filter.PaginationWithInputFilter, totalRows int64) ShiftResponsePagination {
	return ShiftResponsePagination{
		Data: s,
		Meta: p.ToMeta(totalRows),
	}
}
