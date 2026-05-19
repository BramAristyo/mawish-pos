package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedBundlingData(db *gorm.DB) {
	pkgID := uuid.MustParse("00000000-0000-0000-0000-000000000701")
	desc := "Perfect morning combo"

	bundlingPackage := domain.BundlingPackage{
		ID:          pkgID,
		Name:        "Breakfast Bundle",
		Description: &desc,
		Price:       decimal.NewFromInt(35000),
		Cogs:        decimal.NewFromInt(15000),
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&bundlingPackage)

	items := []domain.BundlingItem{
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000711"),
			BundlingPackageID: pkgID,
			ProductID:         uuid.MustParse("00000000-0000-0000-0000-000000000601"), // Espresso
			Qty:               1,
		},
		{
			ID:                uuid.MustParse("00000000-0000-0000-0000-000000000712"),
			BundlingPackageID: pkgID,
			ProductID:         uuid.MustParse("00000000-0000-0000-0000-000000000604"), // Croissant
			Qty:               1,
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&items)
}
