package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaydee029/SeeALie/user/internal/database"
	"github.com/jaydee029/SeeALie/user/middleware"
	"github.com/joho/godotenv"
)

type apiconfig struct {
	jwtsecret string
	DB        *database.Queries
	DBpool    *pgxpool.Pool
}

func main() {
	godotenv.Load(".env")

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

	dbcon, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbcon.Close()

	queries := database.New(dbcon)

	apicfg := &apiconfig{
		jwtsecret: jwtsecret,
		DB:        queries,
		DBpool:    dbcon,
	}

	r := chi.NewRouter()
	s := chi.NewRouter()

	s.Post("/signup", apicfg.Signup)
	s.Post("/login", apicfg.Login)
	s.Post("/refresh", apicfg.Refresh)
	s.Post("/revoke", apicfg.Revoke)
	r.Mount("/auth", s)

	sermux := middleware.Corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The authentication server is live on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
