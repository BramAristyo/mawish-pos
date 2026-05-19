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

type ModifierOptionRepository struct {
	DB *gorm.DB
}

func NewModifierOptionRepository(db *gorm.DB) *ModifierOptionRepository {
	return &ModifierOptionRepository{DB: db}
}

func (r *ModifierOptionRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.ModifierOption, error) {
	var mo []domain.ModifierOption
	var totalRows int64

	allowedFields := map[string]string{
		"name":             "name",
		"price_adjustment": "price_adjustment",
		"cogs_adjustment":  "cogs_adjustment",
		"created_at":       "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.ModifierOption{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, []domain.ModifierOption{}, err
	}

	if err := q.Preload("ModifierGroup").Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&mo).Error; err != nil {
		return 0, []domain.ModifierOption{}, err
	}

	return totalRows, mo, nil
}

func (r *ModifierOptionRepository) FindById(ctx context.Context, id uuid.UUID) (domain.ModifierOption, error) {
	var mo domain.ModifierOption

	if err := r.DB.WithContext(ctx).Preload("ModifierGroup").First(&mo, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ModifierOption{}, usecase_errors.NotFound
		}
		return domain.ModifierOption{}, err
	}

	return mo, nil
}

func (r *ModifierOptionRepository) Store(ctx context.Context, mo *domain.ModifierOption) (domain.ModifierOption, error) {
	if err := r.DB.WithContext(ctx).Create(mo).Error; err != nil {
		return domain.ModifierOption{}, err
	}

	return *mo, nil
}

func (r *ModifierOptionRepository) Update(ctx context.Context, id uuid.UUID, mo *domain.ModifierOption) (domain.ModifierOption, error) {
	var existing domain.ModifierOption
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ModifierOption{}, usecase_errors.NotFound
		}
		return domain.ModifierOption{}, err
	}

	updateData := map[string]any{
		"modifier_group_id": mo.ModifierGroupID,
		"name":              mo.Name,
		"price_adjustment":  mo.PriceAdjustment,
		"cogs_adjustment":   mo.CogsAdjustment,
	}

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.ModifierOption{}, err
	}

	return existing, nil
}

func (r *ModifierOptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.ModifierOption{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *ModifierOptionRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.ModifierOption{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}
