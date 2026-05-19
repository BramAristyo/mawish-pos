package usecase

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
)

type AttendanceUseCase struct {
	Repo      *repository.AttendanceRepository
	ShiftRepo *repository.ShiftScheduleRepository
}

func NewAttendanceUseCase(repo *repository.AttendanceRepository, shiftRepo *repository.ShiftScheduleRepository) *AttendanceUseCase {
	return &AttendanceUseCase{
		Repo:      repo,
		ShiftRepo: shiftRepo,
	}
}

func (u *AttendanceUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.AttendanceResponsePagination, error) {
	totalRows, attendances, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.AttendanceResponsePagination{}, err
	}

	res := dto.ToAttendanceResponses(attendances)
	return dto.ToAttendanceResponsePagination(res, req, totalRows), nil
}

func (u *AttendanceUseCase) Store(ctx context.Context, req dto.AttendanceRequest) (dto.AttendanceResponse, error) {
	attendance, err := dto.ToAttendanceDomain(req)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	if attendance.ShiftScheduleID != nil {
		shift, err := u.ShiftRepo.FindById(ctx, *attendance.ShiftScheduleID)
		if err == nil {
			attendance.CalculateLateness(shift)
		}
	}

	res, err := u.Repo.Store(ctx, &attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	// We might want to preload for the response
	// But for now let's just return the created one
	// Usually Store returns the domain which might not have Preloads
	// Let's re-fetch if needed or just use what we have

	return dto.ToAttendanceResponses([]domain.Attendance{res})[0], nil
}
