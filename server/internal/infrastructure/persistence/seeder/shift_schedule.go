package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedShiftScheduleData(db *gorm.DB) {
	schedules := []domain.ShiftSchedule{
		{
			ID:                  uuid.MustParse("00000000-0000-0000-0000-000000000901"),
			Name:                "Morning Shift (Barista)",
			StartTime:           "07:00:00",
			EndTime:             "15:00:00",
			ToleranceMinutes:    15,
			LateIntervalMinutes: 10,
			LateDeductionAmount: 10000,
		},
		{
			ID:                  uuid.MustParse("00000000-0000-0000-0000-000000000902"),
			Name:                "Evening Shift (Closing)",
			StartTime:           "14:00:00",
			EndTime:             "22:00:00",
			ToleranceMinutes:    15,
			LateIntervalMinutes: 10,
			LateDeductionAmount: 10000,
		},
	}
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&schedules)
}
