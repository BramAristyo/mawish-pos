package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ModifierGroupResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	IsRequired bool      `json:"isRequired"`
	DeletedAt  *string   `json:"deletedAt,omitempty"`
	CreatedAt  string    `json:"createdAt"`
	UpdatedAt  string    `json:"updatedAt"`
}

type ModifierGroupDetailResponse struct {
	ID               uuid.UUID                `json:"id"`
	Name             string                   `json:"name"`
	IsRequired       bool                     `json:"isRequired"`
	DeletedAt        *string                  `json:"deletedAt,omitempty"`
	CreatedAt        string                   `json:"createdAt"`
	UpdatedAt        string                   `json:"updatedAt"`
	Options          []ModifierOptionResponse `json:"options,omitempty"`
	ProductModifiers []ProductResponse        `json:"productModifiers"`
}

type ModifierOptionResponse struct {
	ID              uuid.UUID       `json:"id"`
	ModifierGroupID uuid.UUID       `json:"modifierGroupId"`
	Name            string          `json:"name"`
	PriceAdjustment decimal.Decimal `json:"priceAdjustment"`
	CogsAdjustment  decimal.Decimal `json:"cogsAdjustment"`
	DeletedAt       *string         `json:"deletedAt,omitempty"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`
}

type ModifierGroupResponsePagination struct {
	Data []ModifierGroupResponse `json:"data"`
	Meta filter.Meta             `json:"meta"`
}

type CreateModifierOptionRequest struct {
	Name            string           `json:"name" binding:"required,min=3,max=100"`
	PriceAdjustment *decimal.Decimal `json:"priceAdjustment" binding:"required"`
	CogsAdjustment  *decimal.Decimal `json:"cogsAdjustment" binding:"required"`
}

type UpdateModifierOptionRequest struct {
	ID              *uuid.UUID       `json:"id"`
	Name            string           `json:"name" binding:"required,min=3,max=100"`
	PriceAdjustment *decimal.Decimal `json:"priceAdjustment" binding:"required"`
	CogsAdjustment  *decimal.Decimal `json:"cogsAdjustment" binding:"required"`
}

type CreateModifierGroupRequest struct {
	Name             string                        `json:"name" binding:"required,min=3,max=100"`
	IsRequired       bool                          `json:"isRequired"`
	Options          []CreateModifierOptionRequest `json:"options" binding:"required,min=1,dive"`
	ProductModifiers []uuid.UUID                   `json:"productModifiers"`
}

type UpdateModifierGroupRequest struct {
	Name             string                        `json:"name" binding:"required,min=3,max=100"`
	IsRequired       bool                          `json:"isRequired"`
	Options          []UpdateModifierOptionRequest `json:"options" binding:"required,min=1,dive"`
	ProductModifiers []uuid.UUID                   `json:"productModifiers"`
}

func toModifierOptionResponse(mo *domain.ModifierOption) ModifierOptionResponse {
	resp := ModifierOptionResponse{
		ID:              mo.ID,
		ModifierGroupID: mo.ModifierGroupID,
		Name:            mo.Name,
		PriceAdjustment: mo.PriceAdjustment,
		CogsAdjustment:  mo.CogsAdjustment,
		CreatedAt:       mo.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       mo.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if mo.DeletedAt.Valid {
		at := mo.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToModifierGroupResponse(mg *domain.ModifierGroup) ModifierGroupResponse {
	resp := ModifierGroupResponse{
		ID:         mg.ID,
		Name:       mg.Name,
		IsRequired: mg.IsRequired,
		CreatedAt:  mg.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  mg.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if mg.DeletedAt.Valid {
		at := mg.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToModifierGroupDetailResponse(mg *domain.ModifierGroup) ModifierGroupDetailResponse {
	var options []ModifierOptionResponse
	if len(mg.ModifierOptions) > 0 {
		for _, o := range mg.ModifierOptions {
			options = append(options, toModifierOptionResponse(&o))
		}
	}

	var productModifiers []ProductResponse
	if len(mg.ProductModifiers) > 0 {
		for _, pm := range mg.ProductModifiers {
			productModifiers = append(productModifiers, ToProductResponse(&pm.Product))
		}
	}

	resp := ModifierGroupDetailResponse{
		ID:               mg.ID,
		Name:             mg.Name,
		IsRequired:       mg.IsRequired,
		CreatedAt:        mg.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        mg.UpdatedAt.Format("2006-01-02 15:04:05"),
		Options:          options,
		ProductModifiers: productModifiers,
	}

	if mg.DeletedAt.Valid {
		at := mg.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToModifierGroupResponsePagination(mg []ModifierGroupResponse, f filter.PaginationWithInputFilter, totalRows int64) ModifierGroupResponsePagination {
	return ModifierGroupResponsePagination{
		Data: mg,
		Meta: f.ToMeta(totalRows),
	}
}

func ToModifierGroupModel(req *CreateModifierGroupRequest) domain.ModifierGroup {
	options := make([]domain.ModifierOption, 0, len(req.Options))
	for _, o := range req.Options {
		options = append(options, domain.ModifierOption{
			Name:            o.Name,
			PriceAdjustment: *o.PriceAdjustment,
			CogsAdjustment:  *o.CogsAdjustment,
		})
	}

	productModifiers := make([]domain.ProductModifier, 0, len(req.ProductModifiers))
	for _, pid := range req.ProductModifiers {
		productModifiers = append(productModifiers, domain.ProductModifier{
			ProductID: pid,
		})
	}

	return domain.ModifierGroup{
		Name:             req.Name,
		IsRequired:       req.IsRequired,
		ModifierOptions:  options,
		ProductModifiers: productModifiers,
	}
}

func ToUpdateModifierGroupModel(req *UpdateModifierGroupRequest) domain.ModifierGroup {
	options := make([]domain.ModifierOption, 0, len(req.Options))
	for _, o := range req.Options {
		var id uuid.UUID
		if o.ID != nil {
			id = *o.ID
		}
		options = append(options, domain.ModifierOption{
			ID:              id,
			Name:            o.Name,
			PriceAdjustment: *o.PriceAdjustment,
			CogsAdjustment:  *o.CogsAdjustment,
		})
	}

	productModifiers := make([]domain.ProductModifier, 0, len(req.ProductModifiers))
	for _, pid := range req.ProductModifiers {
		productModifiers = append(productModifiers, domain.ProductModifier{
			ProductID: pid,
		})
	}

	return domain.ModifierGroup{
		Name:             req.Name,
		IsRequired:       req.IsRequired,
		ModifierOptions:  options,
		ProductModifiers: productModifiers,
	}
}
