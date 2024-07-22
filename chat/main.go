package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

	wsServer := newWebServer()
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(wsServer, w, r)
	})

	srv := &http.Server{
		Addr: ":" + port,
		//Handler: sermux,
	}
	log.Printf("The chat server is live on port %s", port)
	log.Fatal(srv.ListenAndServe())

}
