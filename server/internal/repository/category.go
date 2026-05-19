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

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category

	if err := r.DB.WithContext(ctx).Order("updated_at DESC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Category, error) {
	var c []domain.Category
	var totalRows int64

	allowedFields := map[string]string{
		"name":       "name",
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.Model(&domain.Category{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, []domain.Category{}, err
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&c).Error; err != nil {
		return 0, []domain.Category{}, err
	}

	return totalRows, c, nil
}

func (r *CategoryRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	var c domain.Category

	if err := r.DB.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Category{}, usecase_errors.NotFound
		}
		return domain.Category{}, err
	}

	return c, nil
}

func (r *CategoryRepository) Store(ctx context.Context, c *domain.Category) (domain.Category, error) {
	if err := r.DB.WithContext(ctx).Create(c).Error; err != nil {
		return domain.Category{}, err
	}

	return *c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id uuid.UUID, c *domain.Category) (domain.Category, error) {
	var existing domain.Category
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Category{}, usecase_errors.NotFound
		}
		return domain.Category{}, err
	}

	updateData := map[string]any{
		"name":        c.Name,
		"description": c.Description,
	}

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.Category{}, err
	}

	return existing, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.Category{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *CategoryRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.Category{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}
