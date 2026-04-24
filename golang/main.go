package main

import (
	"context"
	"database/sql"
	"errors"
	"golang/controller"
	"golang/internal/config"
	"golang/middleware"
	"golang/migrations"
	"golang/repository"
	"golang/service"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	var handler slog.Handler

	handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	logger := slog.New(handler)

	slog.SetDefault(logger)

	db, err := sql.Open("postgres", cfg.Postgres.ConnString())

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Warn("Failed to init Postgres: ", "error", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.DBTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Warn("Postgres is not ready: ", "error", err)
	} else {
		slog.Info("Postgres is ready")
	}

	err = migrations.RunMigrations(db)

	if err != nil {
		slog.Warn("Failed to init migration process: ", "error", err)
	} else {
		slog.Info("Successfully init migration process")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           0,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	lockManager := repository.NewRedisLock(redisClient)

	ctx, cancel = context.WithTimeout(context.Background(), cfg.Server.RedisTimeout)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		slog.Warn("WARNING: Redis is not reachable, working without cache.", "error", err)
	} else {
		slog.Info("Successfully connected to Redis")
	}

	pgRepo := repository.NewPostgresEventRepository(db)

	repo := repository.NewCachedEventRepository(pgRepo, redisClient, cfg.Server.TTL)

	serv := service.NewEventService(repo, lockManager)

	ctrl := controller.NewEventController(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /events", ctrl.Create)
	mux.HandleFunc("GET /events/{id}", ctrl.Get)
	mux.HandleFunc("GET /events", ctrl.GetAll)
	mux.HandleFunc("PUT /events/{id}", ctrl.Update)
	mux.HandleFunc("DELETE /events/{id}", ctrl.Delete)
	mux.Handle("GET /metrics", promhttp.Handler())

	limiter := middleware.RateLimiter(redisClient, 10, time.Second)

	wrappedMux := limiter(middleware.Logging(mux))

	srv := http.Server{Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), Handler: wrappedMux}

	go func() {
		slog.Info("Server is starting on port", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Warn("Listen error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	slog.Info("Shutdown signal received, shutting down gracefully...")

	ctx, cancel = context.WithTimeout(context.Background(), cfg.Server.TTL)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Warn("Server forced to shutdown: ", "error", err)
	}

	slog.Info("Closing database connections...")
	db.Close()

	slog.Info("Server exited properly")
}
