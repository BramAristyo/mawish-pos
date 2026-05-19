package usecase

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
)

type ShiftScheduleUseCase struct {
	repo *repository.ShiftScheduleRepository
}

func NewShiftScheduleUseCase(repo *repository.ShiftScheduleRepository) *ShiftScheduleUseCase {
	return &ShiftScheduleUseCase{repo: repo}
}

func (u *ShiftScheduleUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.ShiftScheduleResponsePagination, error) {
	totalRows, ss, err := u.repo.Paginate(ctx, req)
	if err != nil {
		return dto.ShiftScheduleResponsePagination{}, err
	}

	responses := dto.ToShiftScheduleResponses(ss)
	return dto.ToShiftScheduleResponsePagination(responses, req, totalRows), nil
}

func (u *ShiftScheduleUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.ShiftScheduleResponse, error) {
	s, err := u.repo.FindById(ctx, id)
	if err != nil {
		return dto.ShiftScheduleResponse{}, err
	}

	return dto.ToShiftScheduleResponse(s), nil
}

func (u *ShiftScheduleUseCase) Store(ctx context.Context, req dto.ShiftScheduleRequest) (dto.ShiftScheduleResponse, error) {
	s := dto.ToShiftScheduleDomain(req)

	result, err := u.repo.Store(ctx, &s)
	if err != nil {
		return dto.ShiftScheduleResponse{}, err
	}

	return dto.ToShiftScheduleResponse(result), nil
}

func (u *ShiftScheduleUseCase) Update(ctx context.Context, id uuid.UUID, req dto.ShiftScheduleRequest) (dto.ShiftScheduleResponse, error) {
	s := dto.ToShiftScheduleDomain(req)
	result, err := u.repo.Update(ctx, id, &s)
	if err != nil {
		return dto.ShiftScheduleResponse{}, err
	}

	return dto.ToShiftScheduleResponse(result), nil
}

func (u *ShiftScheduleUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *ShiftScheduleUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.ShiftScheduleResponse, error) {
	if err := u.repo.Restore(ctx, id); err != nil {
		return dto.ShiftScheduleResponse{}, err
	}

	s, err := u.repo.FindById(ctx, id)
	if err != nil {
		return dto.ShiftScheduleResponse{}, err
	}

	return dto.ToShiftScheduleResponse(s), nil
}

func (u *ShiftScheduleUseCase) GetAll(ctx context.Context) ([]dto.ShiftScheduleResponse, error) {
	ss, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToShiftScheduleResponses(ss), nil
}
