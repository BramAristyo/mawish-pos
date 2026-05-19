package seeder

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedUserData(db *gorm.DB) {
	admin1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	admin2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	cashier1ID := uuid.MustParse("00000000-0000-0000-0000-000000000003")

	users := []*domain.User{
		{
			ID:       admin1ID,
			Name:     "Camelia",
			Email:    "camelia@cameliawhite.my.id",
			Password: "password",
			Role:     domain.RoleAdmin,
		},
		{
			ID:       admin2ID,
			Name:     "White",
			Email:    "white@cameliawhite.my.id",
			Password: "password",
			Role:     domain.RoleAdmin,
		},
		{
			ID:       cashier1ID,
			Name:     "John",
			Email:    "john@cameliawhite.my.id",
			Password: "password",
			Role:     domain.RoleCashier,
		},
	}

	for _, user := range users {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password = string(hashed)
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users)
}
