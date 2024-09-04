package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jaydee029/SeeALie/chat/internal/database"
	"github.com/jaydee029/SeeALie/chat/middleware"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Wserver struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	ChatRooms  map[uuid.UUID]*chatRooms
	jwtSecret  string
	DB         *database.Queries
	Cache      *redis.Client
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("error setting up the redis database" + err.Error())
	}

	Wsserver := &Wserver{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		ChatRooms:  make(map[uuid.UUID]*chatRooms),
		jwtSecret:  jwtsecret,
		DB:         queries,
		Cache:      redisClient,
	}

	r := chi.NewRouter()

	r.Get("/chat", Wsserver.handleChat)
	r.Get("/chat/friends", Wsserver.Getfriends)
	//r.Get("/chat/addfriend", Wsserver.Addfriend)

	sermux := middleware.Corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	//http.HandleFunc("/chat", Wsserver.handleChat)
	log.Printf("The chat server is live on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
