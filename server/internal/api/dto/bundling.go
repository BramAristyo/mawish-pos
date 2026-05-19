package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BundlingItemResponse struct {
	ID                uuid.UUID       `json:"id"`
	BundlingPackageID uuid.UUID       `json:"bundlingPackageId"`
	ProductID         uuid.UUID       `json:"productId"`
	Product           ProductResponse `json:"product"`
	Qty               int             `json:"qty"`
	CreatedAt         string          `json:"createdAt"`
}

type BundlingPackageResponse struct {
	ID            uuid.UUID              `json:"id"`
	Name          string                 `json:"name"`
	Description   *string                `json:"description,omitempty"`
	Price         decimal.Decimal        `json:"price"`
	Cogs          decimal.Decimal        `json:"cogs"`
	ImageURL      *string                `json:"imageUrl"`
	DeletedAt     *string                `json:"deletedAt,omitempty"`
	UpdatedAt     string                 `json:"updatedAt"`
	CreatedAt     string                 `json:"createdAt"`
	BundlingItems []BundlingItemResponse `json:"items,omitempty"`
}

type BundlingPackagePaginationResponse struct {
	Data []BundlingPackageResponse `json:"data"`
	Meta filter.Meta               `json:"meta"`
}

type BundlingItemRequest struct {
	BundlingPackageID *uuid.UUID `json:"bundlingPackageId" binding:"omitempty,uuid"`
	ProductID         uuid.UUID  `json:"productId" binding:"required,uuid"`
	Qty               int        `json:"qty" binding:"required,min=1"`
}

type CreateBundlingPackageRequest struct {
	Name          string                `json:"name" binding:"required,min=3,max=100"`
	Description   *string               `json:"description" binding:"omitempty,max=255"`
	Price         decimal.Decimal       `json:"price" binding:"required,gt=0"`
	Cogs          decimal.Decimal       `json:"cogs" binding:"required,gt=0"`
	ImageURL      *string               `json:"imageUrl" binding:"omitempty,url"`
	BundlingItems []BundlingItemRequest `json:"bundlingItems" binding:"required,min=1,dive"`
}

type UpdateBundlingPackageRequest struct {
	Name          string                `json:"name" binding:"required,min=3,max=100"`
	Description   *string               `json:"description" binding:"omitempty,max=255"`
	Price         decimal.Decimal       `json:"price" binding:"required,gt=0"`
	Cogs          decimal.Decimal       `json:"cogs" binding:"required,gt=0"`
	ImageURL      *string               `json:"imageUrl" binding:"omitempty,url"`
	BundlingItems []BundlingItemRequest `json:"bundlingItems" binding:"required,min=1,dive"`
}

func toBundlingItemResponses(bis []domain.BundlingItem) []BundlingItemResponse {
	bisRes := make([]BundlingItemResponse, len(bis))
	for i, bi := range bis {
		bisRes[i] = BundlingItemResponse{
			ID:                bi.ID,
			BundlingPackageID: bi.BundlingPackageID,
			ProductID:         bi.ProductID,
			Product:           ToProductResponse(&bi.Product),
			Qty:               bi.Qty,
			CreatedAt:         bi.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return bisRes
}

func ToBundlingPackageResponse(bp *domain.BundlingPackage) BundlingPackageResponse {
	resp := BundlingPackageResponse{
		ID:            bp.ID,
		Name:          bp.Name,
		Description:   bp.Description,
		Price:         bp.Price,
		Cogs:          bp.Cogs,
		ImageURL:      bp.ImageURL,
		UpdatedAt:     bp.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:     bp.CreatedAt.Format("2006-01-02 15:04:05"),
		BundlingItems: toBundlingItemResponses(bp.BundlingItems),
	}

	if bp.DeletedAt.Valid {
		at := bp.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func toBundlingPackageResponses(bps []domain.BundlingPackage) []BundlingPackageResponse {
	bpsRes := make([]BundlingPackageResponse, len(bps))
	for i, bp := range bps {
		bpsRes[i] = ToBundlingPackageResponse(&bp)
	}

	return bpsRes
}

func ToBundlingPackagePaginationResponse(bps []domain.BundlingPackage, f filter.PaginationWithInputFilter, totalRows int64) BundlingPackagePaginationResponse {
	return BundlingPackagePaginationResponse{
		Data: toBundlingPackageResponses(bps),
		Meta: f.ToMeta(totalRows),
	}
}

func toBundlingItemModels(req []BundlingItemRequest) []domain.BundlingItem {
	bis := make([]domain.BundlingItem, len(req))
	for i, bi := range req {
		bis[i] = domain.BundlingItem{
			ProductID: bi.ProductID,
			Qty:       bi.Qty,
		}
	}
	return bis
}

func ToBundlingPackageModel(req *CreateBundlingPackageRequest) domain.BundlingPackage {
	return domain.BundlingPackage{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		Cogs:          req.Cogs,
		ImageURL:      req.ImageURL,
		BundlingItems: toBundlingItemModels(req.BundlingItems),
	}
}

func ToUpdateBundlingPackageModel(req *UpdateBundlingPackageRequest) domain.BundlingPackage {
	return domain.BundlingPackage{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		Cogs:          req.Cogs,
		ImageURL:      req.ImageURL,
		BundlingItems: toBundlingItemModels(req.BundlingItems),
	}
}
