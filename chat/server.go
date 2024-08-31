package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jaydee029/SeeALie/chat/internal/auth"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
)

type chatRooms struct {
	ID     uuid.UUID
	Client map[*Client]bool
}

type friendsResponse struct {
	Friends []string
}

func (ws *Wserver) Getfriends(w http.ResponseWriter, r *http.Request) {
	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	userId, err := auth.ValidateToken(token, ws.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	Id, err := uuid.Parse(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	username, err := ws.DB.Get_username(context.Background(), Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding the username:"+err.Error())
		return
	}

	friendsTable, err := ws.DB.Find_friends(context.Background(), username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding friends"+err.Error())
		return
	}

	var friends []string

	for _, i := range friendsTable {
		friends = append(friends, i.Friend)
	}

	respondWithJson(w, http.StatusAccepted, friendsResponse{
		Friends: friends,
	})

}

func (ws *Wserver) handleChat(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	userId, err := auth.ValidateToken(token, ws.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	Id, err := uuid.Parse(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	username, err := ws.DB.Get_username(context.Background(), Id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding the username:"+err.Error())
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	go ws.Runserver()
	targetClient := r.URL.Query().Get("name")

	friendsTable, err := ws.DB.Find_friends(context.Background(), username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding friends"+err.Error())
		return
	}

	var roomid uuid.UUID

	for _, k := range friendsTable {
		if k.Friend == targetClient {
			roomid = k.RoomID
		}
	}

	client := newClient(username, roomid, conn, ws)

	go client.ReadInput()
	go client.WriteInput()

	fmt.Println("New client joined the chat server", client)

	ws.Register <- client

}
