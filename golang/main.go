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
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.Postgres.ConnString())

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to init Postgres: %v", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(); err != nil {
		log.Fatalf("Postgres is not ready: %v", err)
	}

	err = migrations.RunMigrations(db)

	if err != nil {
		log.Fatalf("Failed to init migration process: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	pgRepo := repository.NewPostgresEventRepository(db)

	repo := repository.NewCachedEventRepository(pgRepo, redisClient, cfg.Server.TTL)

	serv := service.NewEventService(repo)

	ctrl := controller.NewEventController(serv)

	mux := http.NewServeMux()

	mux.Handle("/events", middleware.AllowMethod(http.MethodPost)(http.HandlerFunc(ctrl.Create)))

	srv := http.Server{Addr: "localhost:8080", Handler: mux}

	log.Printf("Server starting on port %s", cfg.Server.Port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}
