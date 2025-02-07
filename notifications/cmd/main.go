package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jaydee029/SeeALie/notifications/internal/database"
	service "github.com/jaydee029/SeeALie/notifications/internal/services"
	"github.com/jaydee029/SeeALie/pubsub"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Print("Port not provided")
		port = "8003"
	}

	domain := os.Getenv("DOMAIN")
	adminEmail := os.Getenv("EMAIL")
	adminPasswd := os.Getenv("PASSWD")

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("Database connection string not set")
	}

	dbcon, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbcon.Close()

	queries := database.New(dbcon)

	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	pb := pubsub.New(conn)
	defer conn.Close()

	services := service.NewService(domain, adminEmail, adminPasswd, queries, pb)

	go services.Run(context.Background())

	//http.HandleFunc("/healthz", srv.healthz)

	log.Printf("Starting HTTP server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
