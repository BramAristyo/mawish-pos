package repository

import (
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/persistence/database"
	"gorm.io/gorm"
)

type BaseRepository[TEntity any] struct {
	DB       *gorm.DB
	Preloads database.PreloadEntity
}
