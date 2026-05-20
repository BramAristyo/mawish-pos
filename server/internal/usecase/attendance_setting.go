package usecase

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"gorm.io/gorm"
)

type AttendanceSettingUseCase struct {
	Repo *repository.AttendanceSettingRepository
}

func NewAttendanceSettingUseCase(repo *repository.AttendanceSettingRepository) *AttendanceSettingUseCase {
	return &AttendanceSettingUseCase{
		Repo: repo,
	}
}

func (u *AttendanceSettingUseCase) Get(ctx context.Context) (dto.AttendanceSettingResponse, error) {
	setting, err := u.Repo.Find(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.AttendanceSettingResponse{}, nil
		}
		return dto.AttendanceSettingResponse{}, err
	}

	return dto.ToAttendanceSettingResponse(setting), nil
}

func (u *AttendanceSettingUseCase) Update(ctx context.Context, req dto.AttendanceSettingRequest) (dto.AttendanceSettingResponse, error) {
	setting := dto.ToAttendanceSettingDomain(req)
	updated, err := u.Repo.Upsert(ctx, setting)
	if err != nil {
		return dto.AttendanceSettingResponse{}, err
	}

	return dto.ToAttendanceSettingResponse(updated), nil
}
