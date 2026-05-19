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

type COARepository struct {
	DB *gorm.DB
}

func NewCOARepository(db *gorm.DB) *COARepository {
	return &COARepository{
		DB: db,
	}
}

func (r *COARepository) GetAll(ctx context.Context) ([]domain.ChartOfAccount, error) {
	var coas []domain.ChartOfAccount
	if err := r.DB.WithContext(ctx).Order("name ASC").Find(&coas).Error; err != nil {
		return nil, err
	}
	return coas, nil
}

func (r *COARepository) GetlAllOperational(ctx context.Context) ([]domain.ChartOfAccount, error) {
	var coas []domain.ChartOfAccount
	if err := r.DB.WithContext(ctx).Order("name ASC").Where("is_operational = ?", true).Find(&coas).Error; err != nil {
		return nil, err
	}
	return coas, nil
}

func (r *COARepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.ChartOfAccount, error) {
	var totalRows int64
	coas := make([]domain.ChartOfAccount, 0, req.PaginationInput.PageSize)

	allowedFields := map[string]string{
		"name":       "name",
		"type":       "type",
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.ChartOfAccount{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if totalRows == 0 {
		return 0, nil, nil
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&coas).Error; err != nil {
		return 0, nil, err
	}

	return totalRows, coas, nil
}

func (r *COARepository) FindById(ctx context.Context, id uuid.UUID) (domain.ChartOfAccount, error) {
	var coa domain.ChartOfAccount
	if err := r.DB.WithContext(ctx).First(&coa, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ChartOfAccount{}, usecase_errors.NotFound
		}
		return domain.ChartOfAccount{}, err
	}
	return coa, nil
}

func (r *COARepository) Store(ctx context.Context, coa *domain.ChartOfAccount) (domain.ChartOfAccount, error) {
	if err := r.DB.WithContext(ctx).Create(coa).Error; err != nil {
		return domain.ChartOfAccount{}, err
	}
	return *coa, nil
}

func (r *COARepository) Update(ctx context.Context, id uuid.UUID, coa *domain.ChartOfAccount) (domain.ChartOfAccount, error) {
	var existing domain.ChartOfAccount
	if err := r.DB.WithContext(ctx).First(&existing, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ChartOfAccount{}, usecase_errors.NotFound
		}
		return domain.ChartOfAccount{}, err
	}

	if existing.IsSystem {
		updateData := map[string]any{
			"name": coa.Name,
		}
		if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
			return domain.ChartOfAccount{}, err
		}
	} else {
		updateData := map[string]any{
			"name": coa.Name,
			"type": coa.Type,
		}
		if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
			return domain.ChartOfAccount{}, err
		}
	}

	return r.FindById(ctx, id)
}

func (r *COARepository) Delete(ctx context.Context, id uuid.UUID) error {
	var coa domain.ChartOfAccount
	if err := r.DB.WithContext(ctx).First(&coa, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return usecase_errors.NotFound
		}
		return err
	}

	if coa.IsSystem {
		return usecase_errors.ForbiddenAccess
	}

	return r.DB.WithContext(ctx).Delete(&coa).Error
}

func (r *COARepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.ChartOfAccount{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}
