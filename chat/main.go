package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jaydee029/SeeALie/chat/handler"
	"github.com/jaydee029/SeeALie/chat/internal/database"
	"github.com/joho/godotenv"
)

type Wserver struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	ChatRooms  map[uuid.UUID]*chatRooms
	jwtSecret  string
	DB         *database.Queries
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

	dbcon, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	queries := database.New(dbcon)

	Wsserver := &Wserver{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		ChatRooms:  make(map[uuid.UUID]*chatRooms),
		jwtSecret:  jwtsecret,
		DB:         queries,
	}

	r := chi.NewRouter()

	r.Get("/chat", Wsserver.handleChat)
	r.Get("/chat/friends", Wsserver.Getfriends)

	sermux := handler.Corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	//http.HandleFunc("/chat", Wsserver.handleChat)
	log.Printf("The chat server is live on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
