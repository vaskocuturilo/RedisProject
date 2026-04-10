package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang/controller"
	"golang/migrations"
	"golang/repository"
	"golang/service"
	"golang/utils"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {

	postgresCredential := utils.GetEnvFromFile()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		postgresCredential.Host,
		postgresCredential.Port,
		postgresCredential.UserName,
		postgresCredential.Password,
		postgresCredential.Database)

	db, err := sql.Open("postgres", psqlInfo)

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

	redisConfig := utils.LoadRedisConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       0,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	pgRepo := repository.NewPostgresEventRepository(db)

	repo := repository.NewCachedEventRepository(pgRepo, redisClient, time.Minute*10)

	serv := service.NewEventService(repo)

	ctrl := controller.NewEventController(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("/events", ctrl.Create)

	srv := http.Server{Addr: "localhost:8080", Handler: mux}

	fmt.Printf("Server running at http://localhost:8080\n")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}
