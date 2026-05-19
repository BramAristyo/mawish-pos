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

type ModifierGroupUseCase struct {
	Repo       *repository.ModifierGroupRepository
	LogUseCase *AuditLogUseCase
}

func NewModifierGroupUseCase(repo *repository.ModifierGroupRepository, log *AuditLogUseCase) *ModifierGroupUseCase {
	return &ModifierGroupUseCase{
		Repo:       repo,
		LogUseCase: log,
	}
}

func (u *ModifierGroupUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.ModifierGroupResponsePagination, error) {
	totalRows, mgs, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.ModifierGroupResponsePagination{}, err
	}

	modifierGroups := make([]dto.ModifierGroupResponse, 0, len(mgs))
	for i := range mgs {
		modifierGroups = append(modifierGroups, dto.ToModifierGroupResponse(&mgs[i]))
	}

	return dto.ToModifierGroupResponsePagination(modifierGroups, req, totalRows), nil
}

func (u *ModifierGroupUseCase) GetAll(ctx context.Context) ([]dto.ModifierGroupResponse, error) {
	mgs, err := u.Repo.GetAll(ctx)
	if err != nil {
		return []dto.ModifierGroupResponse{}, err
	}

	res := make([]dto.ModifierGroupResponse, 0, len(mgs))
	for i := range mgs {
		res = append(res, dto.ToModifierGroupResponse(&mgs[i]))
	}

	return res, nil
}

func (u *ModifierGroupUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.ModifierGroupDetailResponse, error) {
	mg, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.ModifierGroupDetailResponse{}, err
	}

	return dto.ToModifierGroupDetailResponse(&mg), nil
}

func (u *ModifierGroupUseCase) Store(ctx context.Context, req dto.CreateModifierGroupRequest) (dto.ModifierGroupDetailResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ModifierGroupDetailResponse{}, err
	}

	mg := dto.ToModifierGroupModel(&req)
	stored, err := u.Repo.Store(ctx, &mg)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.ModifierGroupDetailResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This modifier group name already exists.",
				},
			}
		}
		return dto.ModifierGroupDetailResponse{}, err
	}

	// Fetch detail for response
	detail, err := u.Repo.FindById(ctx, stored.ID)
	if err != nil {
		return dto.ModifierGroupDetailResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityModifierGroup,
		EntityID:    &stored.ID,
		Description: "User created a new modifier group: " + stored.Name,
	})

	return dto.ToModifierGroupDetailResponse(&detail), nil
}

func (u *ModifierGroupUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateModifierGroupRequest) (dto.ModifierGroupDetailResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ModifierGroupDetailResponse{}, err
	}

	mg := dto.ToUpdateModifierGroupModel(&req)
	updated, err := u.Repo.Update(ctx, id, &mg)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.ModifierGroupDetailResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Name",
					Tag:      "unique",
					Value:    req.Name,
					Message:  "This modifier group name already exists.",
				},
			}
		}
		return dto.ModifierGroupDetailResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityModifierGroup,
		EntityID:    &updated.ID,
		Description: "User updated modifier group: " + updated.Name,
	})

	return dto.ToModifierGroupDetailResponse(&updated), nil
}

func (u *ModifierGroupUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return err
	}

	// Fetch before delete to get name for log
	mg, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityModifierGroup,
		EntityID:    &id,
		Description: "User deleted modifier group: " + mg.Name,
	})

	return nil
}

func (u *ModifierGroupUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.ModifierGroupDetailResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ModifierGroupDetailResponse{}, err
	}

	if err := u.Repo.Restore(ctx, id); err != nil {
		return dto.ModifierGroupDetailResponse{}, err
	}

	mg, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.ModifierGroupDetailResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntityModifierGroup,
		EntityID:    &id,
		Description: "User restored modifier group: " + mg.Name,
	})

	return dto.ToModifierGroupDetailResponse(&mg), nil
}
