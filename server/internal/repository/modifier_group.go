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

type ModifierGroupRepository struct {
	DB *gorm.DB
}

func NewModifierGroupRepository(db *gorm.DB) *ModifierGroupRepository {
	return &ModifierGroupRepository{
		DB: db,
	}
}

func (r *ModifierGroupRepository) GetAll(ctx context.Context) ([]domain.ModifierGroup, error) {
	var mg []domain.ModifierGroup
	if err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&mg).Error; err != nil {
		return []domain.ModifierGroup{}, err
	}

	return mg, nil
}

func (r *ModifierGroupRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.ModifierGroup, error) {
	var mg []domain.ModifierGroup
	var totalRows int64

	allowedFields := map[string]string{
		"name":       "name",
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.ModifierGroup{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&mg).Error; err != nil {
		return 0, nil, err
	}

	return totalRows, mg, nil
}

func (r *ModifierGroupRepository) FindById(ctx context.Context, id uuid.UUID) (domain.ModifierGroup, error) {
	var mg domain.ModifierGroup

	err := r.DB.WithContext(ctx).
		Preload("ModifierOptions").
		Preload("ProductModifiers.Product.Category").
		Where("id = ?", id).
		First(&mg).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ModifierGroup{}, usecase_errors.NotFound
		}
		return domain.ModifierGroup{}, err
	}

	return mg, nil
}

func (r *ModifierGroupRepository) Store(ctx context.Context, mg *domain.ModifierGroup) (domain.ModifierGroup, error) {
	if err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(mg).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return domain.ModifierGroup{}, err
	}

	return *mg, nil
}

func (r *ModifierGroupRepository) Update(ctx context.Context, id uuid.UUID, mg *domain.ModifierGroup) (domain.ModifierGroup, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.ModifierGroup
		if err := tx.Where("id = ?", id).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return usecase_errors.NotFound
			}
			return err
		}

		updateData := map[string]any{
			"name":        mg.Name,
			"is_required": mg.IsRequired,
		}

		if err := tx.Model(&existing).Updates(updateData).Error; err != nil {
			return err
		}

		// Sync ModifierOptions
		var currentOptions []domain.ModifierOption
		if err := tx.Where("modifier_group_id = ?", id).Find(&currentOptions).Error; err != nil {
			return err
		}

		currentOptionsMap := make(map[uuid.UUID]bool)
		for _, o := range currentOptions {
			currentOptionsMap[o.ID] = true
		}

		reqOptionsMap := make(map[uuid.UUID]bool)
		for _, o := range mg.ModifierOptions {
			if o.ID != uuid.Nil {
				reqOptionsMap[o.ID] = true
			}
		}

		// Delete options not in request
		for _, o := range currentOptions {
			if !reqOptionsMap[o.ID] {
				if err := tx.Delete(&o).Error; err != nil {
					return err
				}
			}
		}

		// Update or Create options
		for _, o := range mg.ModifierOptions {
			o.ModifierGroupID = id
			if o.ID != uuid.Nil && currentOptionsMap[o.ID] {
				if err := tx.Model(&domain.ModifierOption{}).Where("id = ?", o.ID).Updates(o).Error; err != nil {
					return err
				}
			} else {
				o.ID = uuid.New() // Ensure new ID if not provided or doesn't exist
				if err := tx.Create(&o).Error; err != nil {
					return err
				}
			}
		}

		// Sync ProductModifiers (Join Table)
		var currentProductIds []uuid.UUID
		if err := tx.Model(&domain.ProductModifier{}).Where("modifier_group_id = ?", id).Pluck("product_id", &currentProductIds).Error; err != nil {
			return err
		}

		currentPmMap := make(map[uuid.UUID]bool)
		for _, pid := range currentProductIds {
			currentPmMap[pid] = true
		}

		reqPmMap := make(map[uuid.UUID]bool)
		for _, pm := range mg.ProductModifiers {
			reqPmMap[pm.ProductID] = true
		}

		var pmToDelete []uuid.UUID
		for _, pid := range currentProductIds {
			if !reqPmMap[pid] {
				pmToDelete = append(pmToDelete, pid)
			}
		}

		var pmToCreate []domain.ProductModifier
		for _, pm := range mg.ProductModifiers {
			if !currentPmMap[pm.ProductID] {
				pm.ModifierGroupID = id
				pmToCreate = append(pmToCreate, pm)
			}
		}

		if len(pmToDelete) > 0 {
			if err := tx.Where("modifier_group_id = ? AND product_id IN ?", id, pmToDelete).Delete(&domain.ProductModifier{}).Error; err != nil {
				return err
			}
		}

		if len(pmToCreate) > 0 {
			if err := tx.Create(&pmToCreate).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return domain.ModifierGroup{}, err
	}

	return r.FindById(ctx, id)
}

func (r *ModifierGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.ModifierGroup{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *ModifierGroupRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.ModifierGroup{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}
