package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jaydee029/SeeALie/request/internal/database"
	"github.com/jaydee029/SeeALie/request/protos"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	DB          *database.Queries
	Domain      string
	AdminEmail  string
	AdminPasswd string
	protos.UnimplementedRequestServer
}

func main() {

	godotenv.Load(".env")

	port := os.Getenv("PORT")
	domain := os.Getenv("DOMAIN")
	adminEmail := os.Getenv("EMAIL")
	adminPasswd := os.Getenv("PASSWD")

	if port == "" {
		log.Print("Port not provided")
		port = "8080"
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

	srv := &server{
		DB:          queries,
		Domain:      domain,
		AdminEmail:  adminEmail,
		AdminPasswd: adminPasswd,
	}

	grpcsrv := grpc.NewServer()
	protos.RegisterRequestServer(grpcsrv, srv)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Printf("listener fsiled %v", err)
	}

	log.Printf("The request server is live on port %s", port)
	log.Fatal(grpcsrv.Serve(listener))
}
