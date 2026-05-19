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

// goose -dir internal/infrastructure/persistence/migrations create create_users_table sql

func main() {
	cfg := config.GetConfig()
	zapLogger := logger.NewZapLogger(cfg)

	err := database.InitDb(cfg, zapLogger)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	defer database.CloseDb()

	sqlDb, _ := database.GetDb().DB()

	if len(os.Args) < 2 {
		log.Fatal("usage: migration [up|down|status|reset|version]")
	}

	if os.Args[1] == "reset-seeder" {
		err = goose.RunContext(context.Background(), "reset", sqlDb, "internal/infrastructure/persistence/migrations")
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		err = goose.RunContext(context.Background(), "up", sqlDb, "internal/infrastructure/persistence/migrations")
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		seeder.SeedAll(database.GetDb())
		fmt.Println("Migration with seeder successfully!")
		return
	}

	err = goose.RunContext(context.Background(), os.Args[1], sqlDb, "internal/infrastructure/persistence/migrations")
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration executed successfully!")
}
