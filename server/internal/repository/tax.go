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

type TaxRepository struct {
	DB *gorm.DB
}

func NewTaxRepository(db *gorm.DB) *TaxRepository {
	return &TaxRepository{DB: db}
}

func (r *TaxRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Tax, error) {
	t := make([]domain.Tax, 0, req.PaginationInput.PageSize)
	var totalRows int64

	allowedFields := map[string]string{
		"name":       "name",
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Tax{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, []domain.Tax{}, err
	}

	if totalRows == 0 {
		return 0, []domain.Tax{}, nil
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&t).Error; err != nil {
		return 0, []domain.Tax{}, err
	}

	return totalRows, t, nil
}

func (r *TaxRepository) GetAll(ctx context.Context) ([]domain.Tax, error) {
	var taxes []domain.Tax

	if err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&taxes).Error; err != nil {
		return []domain.Tax{}, err
	}

	return taxes, nil
}

func (r *TaxRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Tax, error) {
	var t domain.Tax

	if err := r.DB.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		return domain.Tax{}, err
	}

	return t, nil
}

func (r *TaxRepository) FindActive(ctx context.Context) (domain.Tax, error) {
	var t domain.Tax

	if err := r.DB.WithContext(ctx).Where("status = ?", true).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Tax{}, usecase_errors.NotFound
		}
		return domain.Tax{}, err
	}

	return t, nil
}

func (r *TaxRepository) Store(ctx context.Context, t *domain.Tax) (domain.Tax, error) {
	if err := r.DB.WithContext(ctx).Create(t).Error; err != nil {
		return domain.Tax{}, err
	}

	return *t, nil
}

func (r *TaxRepository) Update(ctx context.Context, id uuid.UUID, t *domain.Tax) (domain.Tax, error) {
	var existing domain.Tax
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		return domain.Tax{}, err
	}

	updateData := map[string]any{
		"name":       t.Name,
		"percentage": t.Percentage,
	}

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.Tax{}, err
	}

	return existing, nil
}

func (r *TaxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.Tax{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *TaxRepository) Activate(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Deactivate all taxes
		if err := tx.Model(&domain.Tax{}).Where("1 = 1").Update("status", false).Error; err != nil {
			return err
		}

		// Activate the target tax
		result := tx.Model(&domain.Tax{}).Where("id = ?", id).Update("status", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return usecase_errors.NotFound
		}
		return nil
	})
}

func (r *TaxRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Model(&domain.Tax{}).Where("id = ?", id).Update("status", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return nil
}

func (r *TaxRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.Tax{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *TaxRepository) DeleteAll(ctx context.Context) error {
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").Delete(&domain.Tax{}).Error
	if err != nil {
		return err
	}

	return nil
}
