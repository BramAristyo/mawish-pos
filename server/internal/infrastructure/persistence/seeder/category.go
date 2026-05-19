package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedCategoryData(db *gorm.DB) {
	categories := []*domain.Category{
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000101"), Name: "Coffee", Description: "Espresso, brewed, and specialty coffee drinks"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000102"), Name: "Tea", Description: "Hot and iced teas, herbal infusions"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000103"), Name: "Pastries", Description: "Croissants, muffins, scones, and baked goods"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000104"), Name: "Sandwiches", Description: "Freshly made sandwiches and wraps"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000105"), Name: "Salads", Description: "Healthy salads and light meals"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000106"), Name: "Cold Drinks", Description: "Iced coffee, smoothies, and cold beverages"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000107"), Name: "Snacks", Description: "Cookies, chips, and small bites"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000108"), Name: "Breakfast", Description: "Breakfast sets and morning specials"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000109"), Name: "Desserts", Description: "Cakes, brownies, and sweet treats"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000110"), Name: "Merchandise", Description: "Coffee beans, mugs, and branded items"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000111"), Name: "Juice", Description: "Freshly squeezed fruit and vegetable juices"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000112"), Name: "Smoothies", Description: "Blended fruit smoothies and protein shakes"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000113"), Name: "Mocktails", Description: "Non-alcoholic mixed drinks with unique flavors"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000114"), Name: "Hot Drinks", Description: "Hot chocolate, matcha, and warm beverages"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000115"), Name: "Ice Cream", Description: "Various flavors of ice cream and gelato"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000116"), Name: "Waffles", Description: "Sweet and savory waffles"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000117"), Name: "Pancakes", Description: "Fluffy pancakes with various toppings"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000118"), Name: "Rice Bowls", Description: "Rice-based meals with assorted toppings"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000119"), Name: "Pasta", Description: "Italian pasta dishes with rich sauces"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000120"), Name: "Pizza", Description: "Freshly baked pizzas with various toppings"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000121"), Name: "Burgers", Description: "Grilled burgers with beef, chicken, or veggie patties"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000122"), Name: "Fries", Description: "French fries and loaded fries"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000123"), Name: "Steak", Description: "Grilled steak with sides and sauces"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000124"), Name: "Soup", Description: "Warm soups for comfort meals"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000125"), Name: "Vegan", Description: "Plant-based food and beverages"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000126"), Name: "Gluten-Free", Description: "Gluten-free food options"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000127"), Name: "Kids Menu", Description: "Meals specially prepared for kids"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000128"), Name: "Seasonal", Description: "Limited-time seasonal menu items"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000129"), Name: "Signature", Description: "House signature and best-selling items"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000130"), Name: "Combos", Description: "Meal combos and value packages"},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&categories)
}
