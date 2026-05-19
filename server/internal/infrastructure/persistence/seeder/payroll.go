package seeder

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedPayrollData(db *gorm.DB) {
	notes1 := "Regular monthly salary for January 2024"
	notes2 := "Regular monthly salary for January 2024 - Adjusted for late arrivals"

	payrolls := []domain.Payroll{
		{
			ID:             uuid.MustParse("00000000-0000-0000-0000-000000001701"),
			EmployeeID:     uuid.MustParse("00000000-0000-0000-0000-000000000801"), // John Doe
			PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			BaseSalary:     decimal.NewFromInt(5000000),
			TotalDeduction: decimal.NewFromInt(0),
			NetSalary:      decimal.NewFromInt(5000000),
			Notes:          &notes1,
		},
		{
			ID:             uuid.MustParse("00000000-0000-0000-0000-000000001702"),
			EmployeeID:     uuid.MustParse("00000000-0000-0000-0000-000000000801"),
			PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			BaseSalary:     decimal.NewFromInt(4500000),
			TotalDeduction: decimal.NewFromInt(150000),
			NetSalary:      decimal.NewFromInt(4350000),
			Notes:          &notes2,
		},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&payrolls)
}
