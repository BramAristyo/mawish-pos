package usecase

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
)

type CategoryUseCase struct {
	Repo       *repository.CategoryRepository
	LogUseCase *AuditLogUseCase
}

func NewCategoryUseCase(repo *repository.CategoryRepository, log *AuditLogUseCase) *CategoryUseCase {
	return &CategoryUseCase{
		Repo:       repo,
		LogUseCase: log,
	}
}

func (u *CategoryUseCase) GetAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := u.Repo.GetAll(ctx)
	if err != nil {
		return []dto.CategoryResponse{}, nil
	}

	res := make([]dto.CategoryResponse, 0, len(categories))
	for i := range categories {
		res = append(res, dto.ToCategoryResponse(&categories[i]))
	}

	return res, nil
}

func (u *CategoryUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.CategoryResponsePagination, error) {
	totalRows, categories, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.CategoryResponsePagination{}, err
	}

	categoriesResponses := make([]dto.CategoryResponse, 0, len(categories))
	for i := range categories {
		categoriesResponses = append(categoriesResponses, dto.ToCategoryResponse(&categories[i]))
	}

	return dto.ToCategoryResponsePagination(categoriesResponses, req, totalRows), nil
}

func (u *CategoryUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.CategoryResponse, error) {

	category, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	return dto.ToCategoryResponse(&category), nil
}

func (u *CategoryUseCase) Store(ctx context.Context, req dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	category := dto.ToCreateCategoryModel(&req)

	stored, err := u.Repo.Store(ctx, &category)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.CategoryResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This category name already exists.",
				},
			}
		}
		return dto.CategoryResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityCategory,
		EntityID:    &stored.ID,
		Description: "User created a new category: " + stored.Name,
	})

	return dto.ToCategoryResponse(&stored), nil
}

func (u *CategoryUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	category := dto.ToUpdateCategoryModel(&req)
	updated, err := u.Repo.Update(ctx, id, &category)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.CategoryResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This category name already exists.",
				},
			}
		}
		return dto.CategoryResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityCategory,
		EntityID:    &updated.ID,
		Description: "User updated category: " + updated.Name,
	})

	return dto.ToCategoryResponse(&updated), nil
}

func (u *CategoryUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return err
	}

	category, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityCategory,
		EntityID:    &id,
		Description: "User deleted category: " + category.Name,
	})

	return nil
}

func (u *CategoryUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.CategoryResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	if err := u.Repo.Restore(ctx, id); err != nil {
		return dto.CategoryResponse{}, err
	}

	category, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntityCategory,
		EntityID:    &id,
		Description: "User restored category: " + category.Name,
	})

	return dto.ToCategoryResponse(&category), nil
}
