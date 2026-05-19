package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/config"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/persistence/database"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/persistence/seeder"
	"github.com/BramAristyo/mawish-pos/server/pkg/logger"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.GetConfig()
	zapLogger := logger.NewZapLogger(cfg)

	err := database.InitDb(cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal(err.Error())
	}

	defer database.CloseDb()

	db := database.GetDb()
	sqlDb, _ := database.GetDb().DB()

	if len(os.Args) < 2 {
		log.Fatal("usage: seeder [development|production]")
	}

	if os.Args[1] == "development" {
		err = goose.RunContext(context.Background(), "reset", sqlDb, "internal/infrastructure/persistence/migrations")
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		err = goose.RunContext(context.Background(), "up", sqlDb, "internal/infrastructure/persistence/migrations")
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		seeder.SeedAll(db)
		fmt.Println("reset database schema and seed development successfully!")
		return
	}
}
