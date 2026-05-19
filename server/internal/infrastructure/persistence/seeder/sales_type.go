package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedSalesTypeData(db *gorm.DB) {
	dineInID := uuid.MustParse("00000000-0000-0000-0000-000000000501")
	takeAwayID := uuid.MustParse("00000000-0000-0000-0000-000000000502")
	deliveryID := uuid.MustParse("00000000-0000-0000-0000-000000000503")

	salesTypes := []*domain.SalesType{
		{
			ID:   dineInID,
			Name: "Dine In",
			Charges: []domain.AdditionalCharge{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000511"),
					SalesTypeID: dineInID,
					Name:        "Service Charge",
					Type:        domain.Percentage,
					Amount:      decimal.NewFromFloat(5),
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000512"),
					SalesTypeID: dineInID,
					Name:        "PB1",
					Type:        domain.Percentage,
					Amount:      decimal.NewFromFloat(10),
				},
			},
		},
		{
			ID:   takeAwayID,
			Name: "Take Away",
			Charges: []domain.AdditionalCharge{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000521"),
					SalesTypeID: takeAwayID,
					Name:        "Packaging Fee",
					Type:        domain.Fixed,
					Amount:      decimal.NewFromFloat(2000),
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000522"),
					SalesTypeID: takeAwayID,
					Name:        "PB1",
					Type:        domain.Percentage,
					Amount:      decimal.NewFromFloat(10),
				},
			},
		},
		{
			ID:   deliveryID,
			Name: "Delivery",
			Charges: []domain.AdditionalCharge{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000531"),
					SalesTypeID: deliveryID,
					Name:        "Delivery Fee",
					Type:        domain.Fixed,
					Amount:      decimal.NewFromFloat(5000),
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000532"),
					SalesTypeID: deliveryID,
					Name:        "Service Fee",
					Type:        domain.Fixed,
					Amount:      decimal.NewFromFloat(1000),
				},
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000533"),
					SalesTypeID: deliveryID,
					Name:        "PB1",
					Type:        domain.Percentage,
					Amount:      decimal.NewFromFloat(10),
				},
			},
		},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&salesTypes)
}
