package repository

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"gorm.io/gorm"
)

type AttendanceSettingRepository struct {
	DB *gorm.DB
}

func NewAttendanceSettingRepository(db *gorm.DB) *AttendanceSettingRepository {
	return &AttendanceSettingRepository{
		DB: db,
	}
}

func (r *AttendanceSettingRepository) Find(ctx context.Context) (domain.AttendanceSetting, error) {
	var setting domain.AttendanceSetting
	err := r.DB.WithContext(ctx).First(&setting).Error
	return setting, err
}

func (r *AttendanceSettingRepository) Upsert(ctx context.Context, setting domain.AttendanceSetting) (domain.AttendanceSetting, error) {
	var existing domain.AttendanceSetting
	err := r.DB.WithContext(ctx).First(&existing).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := r.DB.WithContext(ctx).Create(&setting).Error; err != nil {
				return domain.AttendanceSetting{}, err
			}
			return setting, nil
		}
		return domain.AttendanceSetting{}, err
	}

	setting.ID = existing.ID
	if err := r.DB.WithContext(ctx).Save(&setting).Error; err != nil {
		return domain.AttendanceSetting{}, err
	}

	return setting, nil
}
