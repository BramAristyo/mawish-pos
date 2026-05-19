package seeder

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedLedgerData(db *gorm.DB) {
	utilsCOA := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	rentCOA := uuid.MustParse("11111111-1111-1111-1111-111111111113")

	note1 := "Electricity Bill - Jan"
	note2 := "Water Bill - Jan"
	note3 := "Store Rent - Jan"

	ledgers := []domain.Ledger{
		{
			ID:              uuid.MustParse("00000000-0000-0000-0000-000000001501"),
			COAID:           utilsCOA,
			Amount:          decimal.NewFromInt(250000),
			Notes:           &note1,
			ReferenceID:     uuid.MustParse("00000000-0000-0000-0000-000000001401"),
			ReferenceType:   domain.LedgerCashTransaction,
			TransactionDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			ID:              uuid.MustParse("00000000-0000-0000-0000-000000001502"),
			COAID:           utilsCOA,
			Amount:          decimal.NewFromInt(100000),
			Notes:           &note2,
			ReferenceID:     uuid.MustParse("00000000-0000-0000-0000-000000001402"),
			ReferenceType:   domain.LedgerCashTransaction,
			TransactionDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.Local),
		},
		{
			ID:              uuid.MustParse("00000000-0000-0000-0000-000000001503"),
			COAID:           rentCOA,
			Amount:          decimal.NewFromInt(5000000),
			Notes:           &note3,
			ReferenceID:     uuid.MustParse("00000000-0000-0000-0000-000000001403"),
			ReferenceType:   domain.LedgerCashTransaction,
			TransactionDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ledgers)
}
