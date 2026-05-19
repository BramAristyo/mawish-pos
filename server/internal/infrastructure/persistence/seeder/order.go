package seeder

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedOrderData(db *gorm.DB) {
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	shiftID := uuid.MustParse("00000000-0000-0000-0000-000000001101")
	salesTypeID := uuid.MustParse("00000000-0000-0000-0000-000000000501")
	taxID := uuid.MustParse("00000000-0000-0000-0000-000000000301")
	prodID := uuid.MustParse("00000000-0000-0000-0000-000000000601")

	order1 := domain.Order{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000001201"),
		ShiftID:     shiftID,
		CashierID:   adminID,
		SalesTypeID: salesTypeID,
		TaxID:       &taxID,
		OrderNumber: "MW/20240101/00001",
		Status:      domain.OrderCompleted,
		Items: []domain.OrderItem{
			{
				ID:           uuid.MustParse("00000000-0000-0000-0000-000000001211"),
				ProductID:    &prodID,
				ProductName:  "Espresso",
				ProductPrice: decimal.NewFromInt(20000),
				ProductCogs:  decimal.NewFromInt(8000),
				Quantity:     2,
			},
		},
		Payments: []domain.Payment{
			{
				ID:     uuid.MustParse("00000000-0000-0000-0000-000000001301"),
				Method: domain.Cash,
				Amount: decimal.NewFromInt(44400), // (20000*2) * 1.11
			},
		},
	}
	order1.CalculateAll()
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&order1)

	voidReason := "Wrong items"
	voidTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local)
	order2 := domain.Order{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000001202"),
		ShiftID:     shiftID,
		CashierID:   adminID,
		SalesTypeID: salesTypeID,
		OrderNumber: "MW/20240101/00002",
		Status:      domain.OrderVoided,
		VoidReason:  &voidReason,
		VoidedBy:    &adminID,
		VoidedAt:    &voidTime,
	}
	order2.CalculateAll()
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&order2)
}
