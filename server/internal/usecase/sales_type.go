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

type SalesTypeUseCase struct {
	Repo       *repository.SalesTypeRepository
	LogUseCase *AuditLogUseCase
}

func NewSalesTypeUseCase(r *repository.SalesTypeRepository, log *AuditLogUseCase) *SalesTypeUseCase {
	return &SalesTypeUseCase{
		Repo:       r,
		LogUseCase: log,
	}
}

func (u *SalesTypeUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.SalesTypeResponsePagination, error) {
	totalRows, salesTypes, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.SalesTypeResponsePagination{}, err
	}

	responses := make([]dto.SalesTypeResponse, 0, len(salesTypes))
	for i := range salesTypes {
		responses = append(responses, dto.ToSalesTypeResponse(&salesTypes[i]))
	}

	return dto.ToSalesTypeResponsePagination(responses, req, totalRows), nil
}

func (u *SalesTypeUseCase) GetAll(ctx context.Context) ([]dto.SalesTypeResponse, error) {
	salesTypes, err := u.Repo.GetAll(ctx)
	if err != nil {
		return []dto.SalesTypeResponse{}, err
	}

	return dto.ToSalesTypeResponses(salesTypes), nil
}

func (u *SalesTypeUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.SalesTypeResponse, error) {
	salesType, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.SalesTypeResponse{}, err
	}

	return dto.ToSalesTypeResponse(&salesType), nil
}

func (u *SalesTypeUseCase) Store(ctx context.Context, req dto.CreateSalesTypeRequest) (dto.SalesTypeResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)
	salesType := dto.ToCreateSalesTypeModel(&req)

	stored, err := u.Repo.Store(ctx, &salesType)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.SalesTypeResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This sales type name already exists.",
				},
			}
		}
		return dto.SalesTypeResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntitySalesType,
		EntityID:    &stored.ID,
		Description: "User created a new sales type: " + stored.Name,
	})

	return dto.ToSalesTypeResponse(&stored), nil
}

func (u *SalesTypeUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateSalesTypeRequest) (dto.SalesTypeResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)
	salesType := dto.ToUpdateSalesTypeModel(&req)

	updated, err := u.Repo.SmartUpdate(ctx, id, &salesType)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.SalesTypeResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This sales type name already exists.",
				},
			}
		}
		return dto.SalesTypeResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntitySalesType,
		EntityID:    &updated.ID,
		Description: "User updated sales type: " + updated.Name,
	})

	return dto.ToSalesTypeResponse(&updated), nil
}

func (u *SalesTypeUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, _ := helper.ExtractUserID(ctx)

	salesType, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntitySalesType,
		EntityID:    &id,
		Description: "User deleted sales type: " + salesType.Name,
	})

	return nil
}

func (u *SalesTypeUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.SalesTypeResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	if err := u.Repo.Restore(ctx, id); err != nil {
		return dto.SalesTypeResponse{}, err
	}

	salesType, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.SalesTypeResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntitySalesType,
		EntityID:    &id,
		Description: "User restored sales type: " + salesType.Name,
	})

	return dto.ToSalesTypeResponse(&salesType), nil
}
