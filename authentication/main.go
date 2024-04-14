package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jaydee029/SeeALie/authentication/handler"
	"github.com/jaydee029/SeeALie/authentication/internal/database"
	"github.com/joho/godotenv"
)

type apiconfig struct {
	jwtsecret string
	DB        *database.Queries
}

func main() {
	godotenv.Load("auth.env")

	port := os.Getenv("PORT")

	if port == "" {
		log.Print("Port not provided")
		port = "8080"
	}

	jwtsecret := os.Getenv("JWT_SECRET")

	if jwtsecret == "" {
		log.Fatalf("JWT Secret not found")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("database connection string not set")
	}

	dbcon, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print(err.Error())
	}
	queries := database.New(dbcon)

	apicfg := &apiconfig{
		jwtsecret: jwtsecret,
		DB:        queries,
	}

	r := chi.NewRouter()
	s := chi.NewRouter()

	s.Post("/signup", apicfg.signup)

	r.Mount("/auth", s)

	sermux := handler.Corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The authentication server is live on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
