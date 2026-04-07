package main

import (
	"database/sql"
	"errors"
	"fmt"
	"golang/controller"
	"golang/repository"
	"golang/service"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	repo := repository.NewPostgresEventRepository(db)

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
