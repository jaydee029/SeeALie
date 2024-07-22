package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  3072,
	WriteBufferSize: 3072,
}

type Client struct {
	conn   *websocket.Conn
	server *wsserver
}

func newClient(server *wsserver, conn *websocket.Conn) *Client {
	return &Client{
		conn:   conn,
		server: server,
	}
}

func serveWS(server *wsserver, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("error:", err)
	}

	client := newClient(server, conn)

	fmt.Println("New client joined")
	fmt.Println(client)
}
