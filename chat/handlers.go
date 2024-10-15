package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	caching "github.com/jaydee029/SeeALie/chat/cache"
	"github.com/jaydee029/SeeALie/chat/internal/auth"
	"github.com/redis/go-redis/v9"
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

func (ws *Wserver) Addfriend(w http.ResponseWriter, r *http.Request) {
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

	username, err := ws.DB.Get_username(r.Context(), Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding the username:"+err.Error())
		return
	}

	targetClient := r.URL.Query().Get("name")
	if targetClient == "" {
		respondWithError(w, http.StatusBadRequest, "no name entered")
	}

	

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

	username, err := ws.DB.Get_username(r.Context(), Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding the username:"+err.Error())
		return
	}

	cacheHit := true
	friendsTable, err := caching.GetCacheFriends(r.Context(), ws.Cache, username)
	if err != nil && err != redis.Nil {
		log.Println("error getting data user profile from redis :  ", err)
		cacheHit = false
	}

	if !cacheHit {
		friendsTable, err = ws.DB.Find_friends(r.Context(), username)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error finding friends"+err.Error())
			return
		}
	}

	err = caching.SetCacheFriends(r.Context(), ws.Cache, username, friendsTable)
	if err != nil {
		log.Println("error setting the cache: ", err)
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

	username, err := ws.DB.Get_username(r.Context(), Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding the username:"+err.Error())
		return
	}

	targetClient := r.URL.Query().Get("name")
	if targetClient == "" {
		respondWithError(w, http.StatusBadRequest, "no name entered")
	}

	cacheHit := true
	friendsTable, err := caching.GetCacheFriends(r.Context(), ws.Cache, username)
	if err != nil && err != redis.Nil {
		log.Println("error getting data user profile from redis :  ", err)
		cacheHit = false
	}

	if !cacheHit {
		friendsTable, err = ws.DB.Find_friends(r.Context(), username)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error finding friends:"+err.Error())
			return
		}
	}

	var roomid uuid.UUID
	friendExists := false
	for _, k := range friendsTable {
		if k.Friend == targetClient {
			roomid = k.RoomID
			friendExists = true
			break
		}
	}
	if !friendExists {
		respondWithError(w, http.StatusInternalServerError, "friend not found:"+err.Error())
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	client := newClient(username, roomid, conn, ws)

	go client.ReadInput()
	go client.WriteInput()

	log.Println("New client joined the chat server", client)

	ws.Register <- client

}
