package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EmployeeID      uuid.UUID `gorm:"type:uuid;not null"`
	Employee        *Employee `gorm:"foreignKey:EmployeeID"`
	Date            time.Time `gorm:"type:date;not null"`
	CheckIn         *time.Time
	CheckOut        *time.Time
	LateMinutes     int     `gorm:"default:0"`
	DeductionAmount float64 `gorm:"type:decimal(12,2);default:0"`
	ShiftScheduleID *uuid.UUID
	ShiftSchedule   *ShiftSchedule `gorm:"foreignKey:ShiftScheduleID"`
	Notes           *string        `gorm:"type:text"`
	PhotoKey        *string        `gorm:"type:varchar(500)"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (a *Attendance) CalculateLateness(policy ShiftSchedule) {
	if a.CheckIn == nil {
		return
	}

	startTime, err := time.Parse("15:04:05", policy.StartTime)
	if err != nil {
		// fallback to without seconds if needed, or just return
		startTime, err = time.Parse("15:04", policy.StartTime)
		if err != nil {
			return
		}
	}

	deadline := time.Date(
		a.Date.Year(), a.Date.Month(), a.Date.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0,
		a.CheckIn.Location(),
	).Add(time.Duration(policy.ToleranceMinutes) * time.Minute)

	if !a.CheckIn.After(deadline) {
		a.LateMinutes = 0
		a.DeductionAmount = 0
		return
	}

	lateMinutes := int(a.CheckIn.Sub(deadline).Minutes())
	intervals := lateMinutes / policy.LateIntervalMinutes

	a.LateMinutes = lateMinutes
	a.DeductionAmount = float64(intervals) * policy.LateDeductionAmount
}
