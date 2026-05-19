package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedDiscountData(db *gorm.DB) {
	discounts := []domain.Discount{
		{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000401"),
			Name:  "Happy Hour 10%",
			Type:  domain.Percentage,
			Value: decimal.NewFromFloat(10),
		},
		{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000402"),
			Name:  "Fixed Discount $5",
			Type:  domain.Fixed,
			Value: decimal.NewFromFloat(5),
		},
		{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000403"),
			Name:  "Weekend Special 15%",
			Type:  domain.Percentage,
			Value: decimal.NewFromFloat(15),
		},
		{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000404"),
			Name:  "Employee Discount 20%",
			Type:  domain.Percentage,
			Value: decimal.NewFromFloat(20),
		},
		{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000405"),
			Name:  "Loyalty Member $10",
			Type:  domain.Fixed,
			Value: decimal.NewFromFloat(10),
		},
		{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000406"),
			Name:  "Student Discount 5%",
			Type:  domain.Percentage,
			Value: decimal.NewFromFloat(5),
		},
		{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000407"),
			Name:  "New Year Promo 25%",
			Type:  domain.Percentage,
			Value: decimal.NewFromFloat(25),
		},
		{
			ID:             uuid.MustParse("00000000-0000-0000-0000-000000000408"),
			Name:           "Happy Custom Discount",
			Type:           domain.Fixed,
			Value:          decimal.NewFromFloat(0.0),
			IsCustomAmount: true,
		},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&discounts)
}
