package main

import (
	"context"
	"database/sql"
	"errors"
	"golang/controller"
	"golang/internal/config"
	"golang/migrations"
	"golang/repository"
	"golang/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	mux.HandleFunc("POST /events", ctrl.Create)
	mux.HandleFunc("GET /events/{id}", ctrl.Get)
	mux.HandleFunc("GET /events", ctrl.GetAll)
	mux.HandleFunc("PUT /events/{id}", ctrl.Update)
	mux.HandleFunc("DELETE /events/{id}", ctrl.Delete)

	srv := http.Server{Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), Handler: mux}

	go func() {
		log.Printf("Server is starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")

	ctx, cancel = context.WithTimeout(context.Background(), cfg.Server.TTL)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Closing database connections...")
	db.Close()

	log.Println("Server exited properly")
}
