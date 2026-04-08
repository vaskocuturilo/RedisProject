package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang/controller"
	"golang/repository"
	"golang/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	databaseUserName := os.Getenv("POSTGRES_USER")
	databasePassword := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, databaseUserName, databasePassword, database)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to init Postgres: %v", err)
		}
	}(db)

	redisAddr := os.Getenv("REDIS_ADDR")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	pgRepo := repository.NewPostgresEventRepository(db)

	repo := repository.NewCachedEventRepository(pgRepo, rdb, time.Minute*10)

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
