package repository

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/persistence/database"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BundlingRepository struct {
	DB *gorm.DB
}

func NewBundlingRepository(db *gorm.DB) *BundlingRepository {
	return &BundlingRepository{DB: db}
}

func (r *BundlingRepository) GetAll(ctx context.Context) ([]domain.BundlingPackage, error) {
	var bps []domain.BundlingPackage
	if err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&bps).Error; err != nil {
		return []domain.BundlingPackage{}, err
	}

	return bps, nil
}

func (r *BundlingRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.BundlingPackage, error) {
	bps := make([]domain.BundlingPackage, 0, req.PaginationInput.PageSize)
	var totalRows int64

	allowedFields := map[string]string{
		"name":       "name",
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.BundlingPackage{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if totalRows == 0 {
		return totalRows, nil, nil
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&bps).Error; err != nil {
		return 0, nil, err
	}

	return totalRows, bps, nil
}

func (r *BundlingRepository) FindById(ctx context.Context, id uuid.UUID) (domain.BundlingPackage, error) {
	var bp domain.BundlingPackage
	if err := r.DB.WithContext(ctx).
		Preload("BundlingItems.Product.Category").
		Where("id = ?", id).
		First(&bp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.BundlingPackage{}, usecase_errors.NotFound
		}
		return domain.BundlingPackage{}, err
	}

	return bp, nil
}

func (r *BundlingRepository) Store(ctx context.Context, bp *domain.BundlingPackage) (domain.BundlingPackage, error) {
	if err := r.DB.WithContext(ctx).Create(bp).Error; err != nil {
		return domain.BundlingPackage{}, err
	}

	return *bp, nil
}

func (r *BundlingRepository) Update(ctx context.Context, id uuid.UUID, bp *domain.BundlingPackage) (domain.BundlingPackage, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.BundlingPackage
		if err := tx.Where("id = ?", id).Preload("BundlingItems").First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return usecase_errors.NotFound
			}
			return err
		}

		updateData := map[string]any{
			"name":        bp.Name,
			"description": bp.Description,
			"price":       bp.Price,
			"cogs":        bp.Cogs,
			"image_url":   bp.ImageURL,
		}

		if err := tx.Model(&existing).Updates(updateData).Error; err != nil {
			return err
		}

		// TODO: Need Smart Sync SOON
		tx.Where("bundling_package_id = ?", id).Delete(&domain.BundlingItem{})
		for i := range bp.BundlingItems {
			bp.BundlingItems[i].BundlingPackageID = id
		}
		tx.Create(&bp.BundlingItems)

		return nil
	})

	if err != nil {
		return domain.BundlingPackage{}, err
	}

	return r.FindById(ctx, id)
}

func (r *BundlingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.BundlingPackage{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *BundlingRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.BundlingPackage{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}
