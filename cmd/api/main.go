package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adi-ch1234/campusslots/internal/config"
	"github.com/adi-ch1234/campusslots/internal/database"
	"github.com/adi-ch1234/campusslots/internal/handler"
	"github.com/adi-ch1234/campusslots/internal/repository"
	"github.com/adi-ch1234/campusslots/internal/service"
	"github.com/adi-ch1234/campusslots/internal/telemetry"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application error %v \n", err)
		os.Exit(1)
	}
}
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Failed to load configuration %v\n", err)

	}
	logger := telemetry.NewLogger(cfg.AppEnv, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	pgPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		return err
	}
	defer pgPool.Close()
	logger.Info("Postgres is connected")
	redisClient, err := database.NewRedisClient(ctx, cfg.RedisAddr)
	if err != nil {
		logger.Warn("redis unavailable at startup; continuing with degraded caching",
			"redis_addr", cfg.RedisAddr,
			"error", err,
		)
	} else {
		logger.Info("connected to redis", "redis_addr", cfg.RedisAddr)
	}
	defer redisClient.Close()
	// 3. Wire repositories -> services -> handlers -> router
	workshopRepo := repository.NewWorkshopRepository(pgPool)
	workshopCache := repository.NewWorkshopCache(redisClient)
	workshopSvc := service.NewWorkshopService(workshopRepo, workshopCache, logger)
	workshopHandler := handler.NewWorkshopHandler(workshopSvc, logger)
	healthHandler := handler.NewHealthHandler(pgPool)
	router := handler.NewRouter(logger, cfg.RequestTimeout, healthHandler, workshopHandler)
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("http server listening", "addr", cfg.HTTPAddr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()
	select {
	case err := <-serverErr:
		return fmt.Errorf("http server failed: %w", err)
	case <-ctx.Done():
		logger.Info("shutdown signal received; draining in-flight requests")
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}
	logger.Info("server stopped cleanly")
	return nil
}
