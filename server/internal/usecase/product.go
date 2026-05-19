package usecase

import (
	"context"
	"fmt"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
)

type ProductUseCase struct {
	Repo       *repository.ProductRepository
	LogUseCase *AuditLogUseCase
}

func NewProductUseCase(repo *repository.ProductRepository, log *AuditLogUseCase) *ProductUseCase {
	return &ProductUseCase{
		Repo:       repo,
		LogUseCase: log,
	}
}

func (u *ProductUseCase) GetAll(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := u.Repo.GetAll(ctx)
	if err != nil {
		return []dto.ProductResponse{}, err
	}

	fmt.Println(len(products))

	res := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		res[i] = dto.ToProductResponse(&product)
	}

	return res, nil
}

func (u *ProductUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.ProductResponsePagination, error) {
	totalRows, products, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.ProductResponsePagination{}, err
	}

	productResponses := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		productResponses = append(productResponses, dto.ToProductResponse(&products[i]))
	}

	return dto.ToProductResponsePagination(productResponses, req, totalRows), nil
}

func (u *ProductUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error) {
	product, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ToProductResponse(&product), nil
}

func (u *ProductUseCase) Store(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	product := dto.ToProductModel(&req)

	stored, err := u.Repo.Store(ctx, &product)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.ProductResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This product name already exists.",
				},
			}
		}
		return dto.ProductResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityProduct,
		EntityID:    &stored.ID,
		Description: "User created a new product: " + stored.Name,
	})

	return dto.ToProductResponse(&stored), nil
}

func (u *ProductUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	product := dto.ToUpdateProductModel(&req)
	updated, err := u.Repo.Update(ctx, id, &product)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.ProductResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This product name already exists.",
				},
			}
		}
		return dto.ProductResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityProduct,
		EntityID:    &updated.ID,
		Description: "User updated product: " + updated.Name,
	})

	return dto.ToProductResponse(&updated), nil
}

func (u *ProductUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return err
	}

	product, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityProduct,
		EntityID:    &id,
		Description: "User deleted product: " + product.Name,
	})

	return nil
}

func (u *ProductUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	if err := u.Repo.Restore(ctx, id); err != nil {
		return dto.ProductResponse{}, err
	}

	product, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntityProduct,
		EntityID:    &id,
		Description: "User restored product: " + product.Name,
	})

	return dto.ToProductResponse(&product), nil
}
