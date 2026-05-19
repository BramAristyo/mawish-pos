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

type TaxUseCase struct {
	Repo       *repository.TaxRepository
	LogUseCase *AuditLogUseCase
}

func NewTaxUseCase(repo *repository.TaxRepository, logUseCase *AuditLogUseCase) *TaxUseCase {
	return &TaxUseCase{
		Repo:       repo,
		LogUseCase: logUseCase,
	}
}

func (u *TaxUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.TaxResponsePagination, error) {
	totalRows, taxes, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.TaxResponsePagination{}, err
	}

	taxResponses := dto.ToTaxResponses(taxes)

	return dto.ToTaxResponsePagination(taxResponses, req, totalRows), nil
}

func (u *TaxUseCase) GetAll(ctx context.Context) ([]dto.TaxResponse, error) {
	taxes, err := u.Repo.GetAll(ctx)
	if err != nil {
		return []dto.TaxResponse{}, err
	}

	return dto.ToTaxResponses(taxes), nil
}

func (u *TaxUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.TaxResponse, error) {
	tax, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.TaxResponse{}, err
	}

	return dto.ToTaxResponse(&tax), nil
}

func (u *TaxUseCase) Store(ctx context.Context, req dto.CreateTaxRequest) (dto.TaxResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)
	tax := dto.ToCreateTaxModel(&req)

	stored, err := u.Repo.Store(ctx, &tax)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.TaxResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This tax name already exists.",
				},
			}
		}
		return dto.TaxResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityTax,
		EntityID:    &stored.ID,
		Description: fmt.Sprintf("Created tax %s with percentage %s", stored.Name, stored.Percentage),
	})

	return dto.ToTaxResponse(&stored), nil
}

func (u *TaxUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTaxRequest) (dto.TaxResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)
	tax := dto.ToUpdateTaxModel(&req)
	updated, err := u.Repo.Update(ctx, id, &tax)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.TaxResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This tax name already exists.",
				},
			}
		}
		return dto.TaxResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityTax,
		EntityID:    &updated.ID,
		Description: fmt.Sprintf("Updated tax %s with percentage %s", updated.Name, updated.Percentage),
	})

	return dto.ToTaxResponse(&updated), nil
}

func (u *TaxUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, _ := helper.ExtractUserID(ctx)

	tax, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityTax,
		EntityID:    &id,
		Description: fmt.Sprintf("Deleted tax %s", tax.Name),
	})

	return nil
}

func (u *TaxUseCase) Activate(ctx context.Context, id uuid.UUID) (dto.TaxResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	if err := u.Repo.Activate(ctx, id); err != nil {
		return dto.TaxResponse{}, err
	}

	tax, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.TaxResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityTax,
		EntityID:    &id,
		Description: fmt.Sprintf("Activated tax %s", tax.Name),
	})

	return dto.ToTaxResponse(&tax), nil
}

func (u *TaxUseCase) Deactivate(ctx context.Context, id uuid.UUID) (dto.TaxResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	if err := u.Repo.Deactivate(ctx, id); err != nil {
		return dto.TaxResponse{}, err
	}

	tax, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.TaxResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityTax,
		EntityID:    &id,
		Description: fmt.Sprintf("Deactivated tax %s", tax.Name),
	})

	return dto.ToTaxResponse(&tax), nil
}

func (u *TaxUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.TaxResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	// Ensure only one tax is active (not deleted)
	if err := u.Repo.DeleteAll(ctx); err != nil {
		return dto.TaxResponse{}, err
	}

	if err := u.Repo.Restore(ctx, id); err != nil {
		return dto.TaxResponse{}, err
	}

	tax, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.TaxResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntityTax,
		EntityID:    &tax.ID,
		Description: fmt.Sprintf("Restored tax %s and set as active", tax.Name),
	})

	return dto.ToTaxResponse(&tax), nil
}
