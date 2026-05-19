package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedEmployeeData(db *gorm.DB) {
	pin1, _ := bcrypt.GenerateFromPassword([]byte("1234"), 10)
	phone1 := "08987654321"

	employees := []domain.Employee{
		{
			ID:         uuid.MustParse("00000000-0000-0000-0000-000000000801"),
			Code:       "EMP-001",
			Phone:      &phone1,
			Name:       "John Doe",
			BaseSalary: 5000000,
			PinHash:    string(pin1),
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&employees)
}
