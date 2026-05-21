package usecase

import (
	"context"
	"fmt"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
)

type AttendanceUseCase struct {
	Repo        *repository.AttendanceRepository
	ShiftRepo   *repository.ShiftScheduleRepository
	StorageRepo domain.StorageRepository
}

func NewAttendanceUseCase(
	repo *repository.AttendanceRepository,
	shiftRepo *repository.ShiftScheduleRepository,
	storageRepo domain.StorageRepository,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		Repo:      repo,
		ShiftRepo: shiftRepo,
		StorageRepo: storageRepo,
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

	key := fmt.Sprintf("attendance/%s/%s/%s.jpg",
		req.EmployeeID,
		req.Date,
		uuid.New().String(),
	)
	uploadUrl, err := u.StorageRepo.GenerateUploadURL(ctx, key)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}
	fmt.Println("URL ", uploadUrl)

	return dto.ToCreateAttedanceResponse(res, uploadUrl), nil
}

func (u *AttendanceUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.AttendanceResponse, error) {
	attendance, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	res := dto.ToAttendanceUpdateResponse(attendance)

	if attendance.PhotoKey != nil {
		getUrl, err := u.StorageRepo.GenerateGetURL(ctx, *attendance.PhotoKey)
		if err != nil {
			return dto.AttendanceResponse{}, err
		}
		res.PhotoURL = getUrl
	}

	return res, nil
}

func (u *AttendanceUseCase) ConfirmAttendanceImage(ctx context.Context, id uuid.UUID, key string) (dto.AttendanceResponse, error) {
	isSuccess, err := u.StorageRepo.VerifyObject(ctx, key)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	if !isSuccess {
		return dto.AttendanceResponse{}, usecase_errors.PhotoVerifyError
	}

	res, err := u.Repo.Update(ctx, id, &domain.Attendance{
		PhotoKey: &key,
	})
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	return dto.ToAttendanceUpdateResponse(res), nil
}

func (u *AttendanceUseCase) Update(ctx context.Context, id uuid.UUID, req dto.AttendanceRequest) (dto.AttendanceResponse, error) {
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

	res, err := u.Repo.Update(ctx, id, &attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	return dto.ToAttendanceUpdateResponse(res), nil
}
