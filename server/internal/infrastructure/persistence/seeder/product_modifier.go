package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedProductModifierData(db *gorm.DB) {
	milkGroupID := uuid.MustParse("00000000-0000-0000-0000-000000000201")

	// Products: Espresso (0601), Latte (0602), Cappuccino (0606)
	mappings := []struct {
		id  string
		pid string
	}{
		{"00000000-0000-0000-0000-000000001601", "00000000-0000-0000-0000-000000000601"},
		{"00000000-0000-0000-0000-000000001602", "00000000-0000-0000-0000-000000000602"},
		{"00000000-0000-0000-0000-000000001603", "00000000-0000-0000-0000-000000000606"},
	}

	var productModifiers []domain.ProductModifier
	for _, m := range mappings {
		productModifiers = append(productModifiers, domain.ProductModifier{
			ID:              uuid.MustParse(m.id),
			ProductID:       uuid.MustParse(m.pid),
			ModifierGroupID: milkGroupID,
		})
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&productModifiers)
}
