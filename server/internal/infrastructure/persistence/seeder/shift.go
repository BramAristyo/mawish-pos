package seeder

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedShiftData(db *gorm.DB) {
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	openAt1 := time.Date(2024, 1, 1, 8, 0, 0, 0, time.Local)
	closeAt1 := time.Date(2024, 1, 1, 16, 0, 0, 0, time.Local)
	openAt2 := time.Date(2024, 1, 1, 16, 0, 0, 0, time.Local)

	notes := "Shift Pagi"
	cash1 := decimal.NewFromInt(1500000)

	shifts := []domain.Shift{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000001101"),
			OpenedBy:    adminID,
			ClosedBy:    &adminID,
			OpeningCash: decimal.NewFromInt(500000),
			ClosingCash: &cash1,
			Notes:       &notes,
			OpenedAt:    openAt1,
			ClosedAt:    &closeAt1,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000001102"),
			OpenedBy:    adminID,
			OpeningCash: decimal.NewFromInt(1500000),
			OpenedAt:    openAt2,
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&shifts)
}
