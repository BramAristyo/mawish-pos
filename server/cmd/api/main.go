package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/api/middleware"
	"github.com/BramAristyo/mawish-pos/server/internal/api/router"
	"github.com/BramAristyo/mawish-pos/server/internal/dependency"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/config"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/persistence/database"
	"github.com/BramAristyo/mawish-pos/server/pkg/logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

var isReady atomic.Bool

func main() {
	// TODO: Implement Dashboard API One API for All use Goroutines! (Business)
	cfg := config.GetConfig()
	zapLogger := logger.NewZapLogger(cfg)

	err := database.InitDb(cfg, zapLogger)
	defer database.CloseDb()

	if err != nil {
		zapLogger.Fatal(err.Error())
	}

	db := database.GetDb()

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	r.Use(ginzap.Ginzap(zapLogger.GetLogger(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zapLogger.GetLogger(), true))

	r.Use(middleware.CORS())
	r.Use(middleware.RateLimiter(zapLogger))

	r.Use(middleware.ErrorHandler(zapLogger))

	handlers := dependency.Bootstrap(db, cfg)

	r.GET("/healthz", func(c *gin.Context) {
		if !isReady.Load() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}

		dbSql, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"reason": "database unreachable",
			})
		}

		if err := dbSql.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"reason": "database unreachable",
			})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	isReady.Store(true)

	r.StaticFile("/health", "./static/health.html")

	// DEPLOY VUE
	r.Static("/assets", "./static/dist/assets")
	r.StaticFile("/favicon.ico", "./static/dist/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/dist/index.html")
	})

	router.RegisterRoutes(r, handlers, cfg)

	s := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	zapLogger.Info("Server existing")
}
