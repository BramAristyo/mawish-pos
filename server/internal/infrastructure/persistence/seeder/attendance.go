package seeder

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAttendanceData(db *gorm.DB) {
	date := time.Date(2026, 5, 20, 0, 0, 0, 0, time.Local)

	// Mock times
	inNormal := time.Date(2026, 5, 20, 06, 55, 0, 0, time.Local) // Early
	outNormal := time.Date(2026, 5, 20, 15, 05, 0, 0, time.Local)
	inLate := time.Date(2026, 5, 20, 14, 35, 0, 0, time.Local) // Late (Tolerance 15m, 14:15 deadline)

	attendances := []domain.Attendance{
		{
			ID:              uuid.MustParse("00000000-0000-0000-0000-000000001001"),
			EmployeeID:      uuid.MustParse("00000000-0000-0000-0000-000000000801"), // John Doe
			Date:            date,
			CheckIn:         &inNormal,
			CheckOut:        &outNormal,
			ShiftScheduleID: uuidPtrSeeder("00000000-0000-0000-0000-000000000901"),
			LateMinutes:     0,
			DeductionAmount: 0,
		},
		{
			ID:              uuid.MustParse("00000000-0000-0000-0000-000000001002"),
			EmployeeID:      uuid.MustParse("00000000-0000-0000-0000-000000000801"), // Jane Smith
			Date:            date,
			CheckIn:         &inLate,
			ShiftScheduleID: uuidPtrSeeder("00000000-0000-0000-0000-000000000902"),
			LateMinutes:     20,
			DeductionAmount: 20000,
		},
	}
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&attendances)
}

func uuidPtrSeeder(s string) *uuid.UUID {
	id := uuid.MustParse(s)
	return &id
}
