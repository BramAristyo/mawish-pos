package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedModifierGroupData(db *gorm.DB) {
	milkID := uuid.MustParse("00000000-0000-0000-0000-000000000201")
	sugarID := uuid.MustParse("00000000-0000-0000-0000-000000000202")
	sizeID := uuid.MustParse("00000000-0000-0000-0000-000000000203")
	iceID := uuid.MustParse("00000000-0000-0000-0000-000000000204")
	toppingID := uuid.MustParse("00000000-0000-0000-0000-000000000205")
	beanID := uuid.MustParse("00000000-0000-0000-0000-000000000206")
	donenessID := uuid.MustParse("00000000-0000-0000-0000-000000000207")
	sideID := uuid.MustParse("00000000-0000-0000-0000-000000000208")

	modifierGroups := []*domain.ModifierGroup{
		{
			ID:         milkID,
			Name:       "Milk Selection",
			IsRequired: false,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002101"), ModifierGroupID: milkID, Name: "Full Cream Milk", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002102"), ModifierGroupID: milkID, Name: "Oat Milk", PriceAdjustment: decimal.NewFromFloat(15000), CogsAdjustment: decimal.NewFromFloat(10000)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002103"), ModifierGroupID: milkID, Name: "Almond Milk", PriceAdjustment: decimal.NewFromFloat(15000), CogsAdjustment: decimal.NewFromFloat(10000)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002104"), ModifierGroupID: milkID, Name: "Soy Milk", PriceAdjustment: decimal.NewFromFloat(12000), CogsAdjustment: decimal.NewFromFloat(8000)},
			},
		},
		{
			ID:         sugarID,
			Name:       "Sugar Level",
			IsRequired: true,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002201"), ModifierGroupID: sugarID, Name: "Normal Sugar", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002202"), ModifierGroupID: sugarID, Name: "Less Sugar (50%)", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002203"), ModifierGroupID: sugarID, Name: "No Sugar", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002204"), ModifierGroupID: sugarID, Name: "Extra Sugar", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
			},
		},
		{
			ID:         sizeID,
			Name:       "Cup Size",
			IsRequired: true,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002301"), ModifierGroupID: sizeID, Name: "Regular", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002302"), ModifierGroupID: sizeID, Name: "Large", PriceAdjustment: decimal.NewFromFloat(5000), CogsAdjustment: decimal.NewFromFloat(2000)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002303"), ModifierGroupID: sizeID, Name: "Extra Large", PriceAdjustment: decimal.NewFromFloat(10000), CogsAdjustment: decimal.NewFromFloat(4000)},
			},
		},
		{
			ID:         iceID,
			Name:       "Ice Level",
			IsRequired: true,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002401"), ModifierGroupID: iceID, Name: "Normal Ice", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002402"), ModifierGroupID: iceID, Name: "Less Ice", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002403"), ModifierGroupID: iceID, Name: "No Ice", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
			},
		},
		{
			ID:         toppingID,
			Name:       "Extra Toppings",
			IsRequired: false,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002501"), ModifierGroupID: toppingID, Name: "Pearl (Boba)", PriceAdjustment: decimal.NewFromFloat(5000), CogsAdjustment: decimal.NewFromFloat(2000)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002502"), ModifierGroupID: toppingID, Name: "Grass Jelly", PriceAdjustment: decimal.NewFromFloat(4000), CogsAdjustment: decimal.NewFromFloat(1500)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002503"), ModifierGroupID: toppingID, Name: "Whipped Cream", PriceAdjustment: decimal.NewFromFloat(6000), CogsAdjustment: decimal.NewFromFloat(2500)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002504"), ModifierGroupID: toppingID, Name: "Caramel Drizzle", PriceAdjustment: decimal.NewFromFloat(3000), CogsAdjustment: decimal.NewFromFloat(1000)},
			},
		},
		{
			ID:         beanID,
			Name:       "Coffee Bean Type",
			IsRequired: true,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002601"), ModifierGroupID: beanID, Name: "House Blend", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002602"), ModifierGroupID: beanID, Name: "100% Arabica", PriceAdjustment: decimal.NewFromFloat(5000), CogsAdjustment: decimal.NewFromFloat(2000)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002603"), ModifierGroupID: beanID, Name: "Single Origin", PriceAdjustment: decimal.NewFromFloat(8000), CogsAdjustment: decimal.NewFromFloat(3000)},
			},
		},
		{
			ID:         donenessID,
			Name:       "Meat Doneness",
			IsRequired: true,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002701"), ModifierGroupID: donenessID, Name: "Rare", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002702"), ModifierGroupID: donenessID, Name: "Medium Rare", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002703"), ModifierGroupID: donenessID, Name: "Medium", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002704"), ModifierGroupID: donenessID, Name: "Well Done", PriceAdjustment: decimal.NewFromFloat(0), CogsAdjustment: decimal.NewFromFloat(0)},
			},
		},
		{
			ID:         sideID,
			Name:       "Side Dish",
			IsRequired: false,
			ModifierOptions: []domain.ModifierOption{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002801"), ModifierGroupID: sideID, Name: "French Fries", PriceAdjustment: decimal.NewFromFloat(15000), CogsAdjustment: decimal.NewFromFloat(5000)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002802"), ModifierGroupID: sideID, Name: "Side Salad", PriceAdjustment: decimal.NewFromFloat(12000), CogsAdjustment: decimal.NewFromFloat(4000)},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000002803"), ModifierGroupID: sideID, Name: "Mashed Potato", PriceAdjustment: decimal.NewFromFloat(18000), CogsAdjustment: decimal.NewFromFloat(6000)},
			},
		},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&modifierGroups)
}
