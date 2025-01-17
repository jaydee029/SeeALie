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
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Print("Port not provided")
		port = "8080"
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

	services := &service.Service{
		Domain:      domain,
		AdminEmail:  adminEmail,
		AdminPasswd: adminPasswd,
		DB:          queries,
	}

	h := &handler.Handler{
		DB:  queries,
		Svc: services,
	}

	go h.Run(context.Background())

	//http.HandleFunc("/healthz", srv.healthz)

	log.Printf("Starting HTTP server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
