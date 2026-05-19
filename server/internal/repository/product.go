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

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product

	if err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&products).Error; err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (r *ProductRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Product, error) {
	p := make([]domain.Product, 0, req.PaginationInput.PageSize)
	var totalRows int64

	allowedFields := map[string]string{
		"name":        "name",
		"category_id": "category_id",
		"price":       "price",
		"cogs":        "cogs",
		"created_at":  "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Product{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if err := q.
		Preload("Category").
		// Preload("ProductModifiers.ModifierGroup.ModifierOptions").
		Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&p).Error; err != nil {
		return 0, nil, err
	}

	return totalRows, p, nil
}

func (r *ProductRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var p domain.Product

	if err := r.DB.WithContext(ctx).
		Preload("Category").
		Preload("ProductModifiers.ModifierGroup.ModifierOptions").
		Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Product{}, usecase_errors.NotFound
		}
		return domain.Product{}, err
	}

	return p, nil
}

func (r *ProductRepository) Store(ctx context.Context, p *domain.Product) (domain.Product, error) {
	if err := r.DB.WithContext(ctx).Create(p).Error; err != nil {
		return domain.Product{}, err
	}

	return *p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id uuid.UUID, p *domain.Product) (domain.Product, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.Product
		if err := tx.Where("id = ?", id).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return usecase_errors.NotFound
			}
			return err
		}

		updateData := map[string]any{
			"category_id": p.CategoryID,
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"cogs":        p.Cogs,
			"image_url":   p.ImageURL,
		}

		if err := tx.Model(&existing).Updates(updateData).Error; err != nil {
			return err
		}

		// Sync ProductModifiers
		var existingIds []uuid.UUID
		tx.Model(&domain.ProductModifier{}).Where("product_id = ?", id).Pluck("modifier_group_id", &existingIds)

		existingMap := make(map[uuid.UUID]bool)
		for _, id := range existingIds {
			existingMap[id] = true
		}

		reqMap := make(map[uuid.UUID]bool)
		for _, pm := range p.ProductModifiers {
			reqMap[pm.ModifierGroupID] = true
		}

		var toDelete []uuid.UUID
		for _, id := range existingIds {
			if !reqMap[id] {
				toDelete = append(toDelete, id)
			}
		}

		var toCreate []domain.ProductModifier
		for _, pm := range p.ProductModifiers {
			if !existingMap[pm.ModifierGroupID] {
				pm.ProductID = id
				toCreate = append(toCreate, pm)
			}
		}

		if len(toDelete) > 0 {
			tx.Where("product_id = ? AND modifier_group_id IN ?", id, toDelete).
				Delete(&domain.ProductModifier{})
		}

		if len(toCreate) > 0 {
			tx.Create(&toCreate)
		}

		return nil
	})

	if err != nil {
		return domain.Product{}, err
	}

	return r.FindById(ctx, id)
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *ProductRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *ProductRepository) FindByCategoryId(ctx context.Context, categoryId uuid.UUID, req filter.PaginationWithInputFilter) (int64, []domain.Product, error) {
	p := make([]domain.Product, 0, req.PaginationInput.PageSize)
	var totalRows int64

	allowedFields := map[string]string{
		"name":       "name",
		"price":      "price",
		"cogs":       "cogs",
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Product{}).Where("category_id = ?", categoryId), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if err := q.
		Preload("Category").
		Preload("ProductModifiers.ModifierGroup.ModifierOptions").
		Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&p).Error; err != nil {
		return 0, nil, err
	}

	return totalRows, p, nil
}
