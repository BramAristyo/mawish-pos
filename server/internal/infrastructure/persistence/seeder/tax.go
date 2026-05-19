package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedTaxData(db *gorm.DB) {
	taxes := []domain.Tax{
		{
			ID:         uuid.MustParse("00000000-0000-0000-0000-000000000301"),
			Name:       "TAX 11%",
			Percentage: decimal.NewFromFloat(11),
		},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&taxes)
}
