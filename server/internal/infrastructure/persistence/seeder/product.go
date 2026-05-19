package seeder

import (
	"fmt"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedProductData(db *gorm.DB) {
	categoryCoffeeID := uuid.MustParse("00000000-0000-0000-0000-000000000101")

	productNames := []string{
		"Espresso", "Latte", "Green Tea", "Croissant", "Turkey Sandwich",
		"Cappuccino", "Americano", "Mocha", "Flat White", "Cold Brew",
		"Earl Grey Tea", "Chamomile Tea", "Iced Peach Tea", "Blueberry Muffin", "Chocolate Chip Cookie",
		"Bagel with Cream Cheese", "Avocado Toast", "Caesar Salad", "Grilled Cheese", "Veggie Wrap",
		"Club Sandwich", "Beef Burger", "Margherita Pizza", "Pasta Carbonara", "Lemonade",
		"Smoothies", "Hot Chocolate", "Scone", "Quiche", "Chicken Caesar Wrap",
	}

	products := make([]domain.Product, len(productNames))
	for i, name := range productNames {
		idStr := fmt.Sprintf("00000000-0000-0000-0000-0000000006%02d", i+1)
		desc := "Description for " + name
		products[i] = domain.Product{
			ID:          uuid.MustParse(idStr),
			Name:        name,
			Description: &desc,
			Price:       decimal.NewFromInt(int64(20000 + (i * 1000))),
			Cogs:        decimal.NewFromInt(int64(8000 + (i * 500))),
			CategoryID:  categoryCoffeeID,
		}
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&products)
}
