package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AdditionalChargeResponse struct {
	ID     uuid.UUID             `json:"id"`
	Name   string                `json:"name"`
	Type   domain.AdjustmentType `json:"type"`
	Amount decimal.Decimal       `json:"amount"`
}

type SalesTypeResponse struct {
	ID        uuid.UUID                  `json:"id"`
	Name      string                     `json:"name"`
	Charges   []AdditionalChargeResponse `json:"charges,omitempty"`
	DeletedAt *string                    `json:"deletedAt,omitempty"`
	CreatedAt string                     `json:"createdAt"`
}

type SalesTypeResponsePagination struct {
	Data []SalesTypeResponse `json:"data"`
	Meta filter.Meta         `json:"meta"`
}

type CreateAdditionalChargeRequest struct {
	Name   string                `json:"name" binding:"required,min=3,max=100"`
	Type   domain.AdjustmentType `json:"type" binding:"required,oneof=percentage fixed"`
	Amount decimal.Decimal       `json:"amount" binding:"required,gt=0"`
}

type CreateSalesTypeRequest struct {
	Name    string                          `json:"name" binding:"required,min=3,max=100"`
	Charges []CreateAdditionalChargeRequest `json:"charges" binding:"omitempty,dive"`
}

type UpdateAdditionalChargeRequest struct {
	ID     *string               `json:"id" binding:"omitempty,uuid"`
	Name   string                `json:"name" binding:"required,min=3,max=100"`
	Type   domain.AdjustmentType `json:"type" binding:"required,oneof=percentage fixed"`
	Amount decimal.Decimal       `json:"amount" binding:"required,gt=0"`
}

type UpdateSalesTypeRequest struct {
	Name    string                          `json:"name" binding:"required,min=3,max=100"`
	Charges []UpdateAdditionalChargeRequest `json:"charges" binding:"omitempty,dive"`
}

func ToAdditionalChargeResponse(c *domain.AdditionalCharge) AdditionalChargeResponse {
	return AdditionalChargeResponse{
		ID:     c.ID,
		Name:   c.Name,
		Type:   c.Type,
		Amount: c.Amount,
	}
}

func ToSalesTypeResponse(s *domain.SalesType) SalesTypeResponse {
	charges := make([]AdditionalChargeResponse, 0, len(s.Charges))
	for _, c := range s.Charges {
		charges = append(charges, ToAdditionalChargeResponse(&c))
	}

	resp := SalesTypeResponse{
		ID:        s.ID,
		Name:      s.Name,
		Charges:   charges,
		CreatedAt: s.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if s.DeletedAt.Valid {
		at := s.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToSalesTypeResponses(s []domain.SalesType) []SalesTypeResponse {
	res := make([]SalesTypeResponse, len(s))
	for i := range s {
		res[i] = ToSalesTypeResponse(&s[i])
	}

	return res
}

func ToSalesTypeResponsePagination(s []SalesTypeResponse, p filter.PaginationWithInputFilter, totalRows int64) SalesTypeResponsePagination {
	return SalesTypeResponsePagination{
		Data: s,
		Meta: p.ToMeta(totalRows),
	}
}

func ToCreateSalesTypeModel(req *CreateSalesTypeRequest) domain.SalesType {
	charges := make([]domain.AdditionalCharge, 0, len(req.Charges))
	for _, c := range req.Charges {
		charges = append(charges, domain.AdditionalCharge{
			Name:   c.Name,
			Type:   c.Type,
			Amount: c.Amount,
		})
	}

	return domain.SalesType{
		Name:    req.Name,
		Charges: charges,
	}
}

func ToUpdateSalesTypeModel(req *UpdateSalesTypeRequest) domain.SalesType {
	charges := make([]domain.AdditionalCharge, 0, len(req.Charges))
	for _, c := range req.Charges {
		var id uuid.UUID

		if c.ID != nil {
			parsedID, err := uuid.Parse(*c.ID)
			if err == nil {
				id = parsedID
			}
		}

		charges = append(charges, domain.AdditionalCharge{
			ID:     id,
			Name:   c.Name,
			Type:   c.Type,
			Amount: c.Amount,
		})
	}

	return domain.SalesType{
		Name:    req.Name,
		Charges: charges,
	}
}
