package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AttendanceUseCase struct {
	Repo                  *repository.AttendanceRepository
	ShiftRepo             *repository.ShiftScheduleRepository
	EmployeeRepo          *repository.EmployeeRepository
	AttendanceSettingRepo *repository.AttendanceSettingRepository
	StorageRepo           domain.StorageRepository
}

func NewAttendanceUseCase(
	repo *repository.AttendanceRepository,
	shiftRepo *repository.ShiftScheduleRepository,
	employeeRepo *repository.EmployeeRepository,
	attendanceSettingRepo *repository.AttendanceSettingRepository,
	storageRepo domain.StorageRepository,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		Repo:                  repo,
		ShiftRepo:             shiftRepo,
		EmployeeRepo:          employeeRepo,
		AttendanceSettingRepo: attendanceSettingRepo,
		StorageRepo:           storageRepo,
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

func (u *AttendanceUseCase) Store(ctx context.Context, req dto.CreateAttendanceRequest) (dto.AttendanceResponse, error) {
	employee, err := u.EmployeeRepo.FindByCode(ctx, req.EmployeeCode)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.PinHash), []byte(req.Pin))
	if err != nil {
		return dto.AttendanceResponse{}, usecase_errors.InvalidPassword
	}

	attendance, err := dto.ToAttendanceDomain(req)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	attendance.EmployeeID = employee.ID

	if attendance.CheckIn != nil {
		schedules, err := u.ShiftRepo.GetAll(ctx)
		if err == nil {
			var activeShift *domain.ShiftSchedule
			checkInTime := attendance.CheckIn.Format("15:04:05")

			for _, s := range schedules {
				start, _ := time.Parse("15:04:05", s.StartTime)
				if s.StartTime == "" || len(s.StartTime) < 5 { // Fallback for 15:04
					start, _ = time.Parse("15:04", s.StartTime)
				}

				end, _ := time.Parse("15:04:05", s.EndTime)
				if s.EndTime == "" || len(s.EndTime) < 5 {
					end, _ = time.Parse("15:04", s.EndTime)
				}

				nowTime, _ := time.Parse("15:04:05", checkInTime)

				bufferStart := start.Add(-2 * time.Hour)
				bufferEnd := end.Add(2 * time.Hour)

				if end.Before(start) {
					if nowTime.After(bufferStart) || nowTime.Before(bufferEnd) {
						activeShift = &s
						break
					}
				} else {
					if nowTime.After(bufferStart) && nowTime.Before(bufferEnd) {
						activeShift = &s
						break
					}
				}
			}

			if activeShift != nil {
				attendance.ShiftScheduleID = &activeShift.ID
				attendance.CalculateLateness(*activeShift)
			}
		}
	}

	if req.Lat != nil && req.Lng != nil {
		setting, err := u.AttendanceSettingRepo.Find(ctx)
		if err == nil {
			if setting.IsValidDistance(*req.Lat, *req.Lng) {
				attendance.LocationStatus = domain.LocationStatusInArea
			} else {
				attendance.LocationStatus = domain.LocationStatusOutArea
			}
		}
	}

	res, err := u.Repo.Store(ctx, &attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	key := fmt.Sprintf("attendance/%s/%s/%s.jpg",
		employee.ID,
		req.Date,
		uuid.New().String(),
	)
	uploadUrl, err := u.StorageRepo.GenerateUploadURL(ctx, key)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

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

func (u *AttendanceUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateAttendanceRequest) (dto.AttendanceResponse, error) {
	attendance, err := dto.ToUpdateAttendanceDomain(req)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	if attendance.ShiftScheduleID != nil {
		shift, err := u.ShiftRepo.FindById(ctx, *attendance.ShiftScheduleID)
		if err == nil {
			attendance.CalculateLateness(shift)
		}
	}

	// Calculate Location Status
	if req.Lat != nil && req.Lng != nil {
		setting, err := u.AttendanceSettingRepo.Find(ctx)
		if err == nil {
			if setting.IsValidDistance(*req.Lat, *req.Lng) {
				attendance.LocationStatus = domain.LocationStatusInArea
			} else {
				attendance.LocationStatus = domain.LocationStatusOutArea
			}
		}
	}

	res, err := u.Repo.Update(ctx, id, &attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	return dto.ToAttendanceUpdateResponse(res), nil
}
