package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func ServeWS(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("error:", err)
	}

	client := newClient(conn)
	fmt.Println("New client joined")
	fmt.Println(client)
}
