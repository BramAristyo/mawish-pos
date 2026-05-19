package usecase

import (
	"context"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
)

type ShiftUseCase struct {
	Repo       *repository.ShiftRepository
	LogUseCase *AuditLogUseCase
}

func NewShiftUseCase(repo *repository.ShiftRepository, log *AuditLogUseCase) *ShiftUseCase {
	return &ShiftUseCase{
		Repo:       repo,
		LogUseCase: log,
	}
}

func (u *ShiftUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.ShiftResponsePagination, error) {
	totalRows, shifts, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.ShiftResponsePagination{}, err
	}

	res := make([]dto.ShiftResponse, 0, len(shifts))
	for i := range shifts {
		res = append(res, dto.ToShiftResponse(&shifts[i]))
	}

	return dto.ToShiftResponsePagination(res, req, totalRows), nil
}

func (u *ShiftUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.ShiftResponse, error) {
	shift, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	return dto.ToShiftResponse(&shift), nil
}

func (u *ShiftUseCase) OpenShift(ctx context.Context, req dto.OpenShiftRequest) (dto.ShiftResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	// Check if user already has an open shift
	existing, _ := u.Repo.FindOpenShiftByUserId(ctx, userId)
	if existing.ID != uuid.Nil {
		return dto.ShiftResponse{}, usecase_errors.ShiftAlreadyOpen
	}

	shift := domain.Shift{
		OpenedBy:    userId,
		OpeningCash: req.OpeningCash,
		Notes:       req.Notes,
		OpenedAt:    time.Now(),
	}

	stored, err := u.Repo.Store(ctx, &shift)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityShift,
		EntityID:    &stored.ID,
		Description: "User opened a new shift",
	})

	return u.FindById(ctx, stored.ID)
}

func (u *ShiftUseCase) CloseShift(ctx context.Context, req dto.CloseShiftRequest) (dto.ShiftResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	shift, err := u.Repo.FindOpenShiftByUserId(ctx, userId)
	if err != nil {
		return dto.ShiftResponse{}, usecase_errors.NoOpenShift
	}

	now := time.Now()
	shift.ClosedBy = &userId
	shift.ClosingCash = &req.ClosingCash
	shift.Notes = req.Notes
	shift.ClosedAt = &now

	updated, err := u.Repo.Update(ctx, shift.ID, &shift)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityShift,
		EntityID:    &updated.ID,
		Description: "User closed the shift",
	})

	return dto.ToShiftResponse(&updated), nil
}

func (u *ShiftUseCase) FindOpenShiftByCurrent(ctx context.Context) (dto.ShiftResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	shift, err := u.Repo.FindOpenShiftByUserId(ctx, userId)
	if err != nil {
		return dto.ShiftResponse{}, usecase_errors.NoOpenShift
	}

	return dto.ToShiftResponse(&shift), nil
}

func (u *ShiftUseCase) UpsertExpenses(ctx context.Context, req dto.UpsertShiftExpensesRequest) (dto.ShiftResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	shift, err := u.Repo.FindOpenShiftByUserId(ctx, userId)
	if err != nil {
		return dto.ShiftResponse{}, usecase_errors.NoOpenShift
	}

	updated, err := u.Repo.Update(ctx, shift.ID, &shift)
	if err != nil {
		return dto.ShiftResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityShift,
		EntityID:    &updated.ID,
		Description: "User updated shift expenses",
	})

	return dto.ToShiftResponse(&updated), nil
}
