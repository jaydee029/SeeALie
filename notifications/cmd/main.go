package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	service "github.com/jaydee029/SeeALie/request"
	"github.com/jaydee029/SeeALie/request/handler"
	"github.com/jaydee029/SeeALie/request/internal/database"
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
	defer conn.Close()

	services := &service.Service{
		Domain:      domain,
		AdminEmail:  adminEmail,
		AdminPasswd: adminPasswd,
		DB:          queries,
		Pubsub:      conn,
	}

	h := &handler.PubSubhandler{
		Svc: services,
	}

	go h.Svc.Run(context.Background())

	//http.HandleFunc("/healthz", srv.healthz)

	log.Printf("Starting HTTP server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
