package seeder

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedCashTransactionData(db *gorm.DB) {
	desc1 := "Buying cleaning supplies (Soap and Mop)"
	desc2 := "Coffee Workshop Registration Fee"

	transactions := []domain.CashTransaction{
		{
			ID:          uuid.MustParse("66666666-1111-1111-1111-111111111111"),
			COAID:       uuid.MustParse("11111111-1111-1111-1111-111111111111"), // Utilities
			Type:        domain.CashOut,
			Amount:      decimal.NewFromInt(55000),
			Description: &desc1,
			Date:        time.Now(),
		},
		{
			ID:          uuid.MustParse("66666666-1111-1111-1111-111111111112"),
			COAID:       uuid.MustParse("11111111-1111-1111-1111-111111111116"), // Other Income
			Type:        domain.CashIn,
			Amount:      decimal.NewFromInt(150000),
			Description: &desc2,
			Date:        time.Now(),
		},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&transactions)
}
