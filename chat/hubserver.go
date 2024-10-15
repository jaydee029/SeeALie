package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func (ws *Wserver) Runserver(ctx context.Context) {

	for {
		select {
		case Client := <-ws.Register:
			ws.registerClient(ctx, Client)
		case Client := <-ws.Unregister:
			ws.unregisterClient(ctx, Client)
		case Message := <-ws.Broadcast:
			ws.BroadcastMessage(ctx, Message)
		}

	}
}

func (ws *Wserver) registerClient(ctx context.Context, client *Client) {
	room, ok := ws.ChatRooms[client.Roomid]
	if ok {
		room.Client[client] = true
		msg := &Message{
			content: []byte(client.username + "has joined the chat"),
			roomid:  client.Roomid,
			sender:  client.username,
		}
		ws.BroadcastMessage(ctx, msg)
	} else {
		ws.ChatRooms[client.Roomid] = &chatRooms{
			ID:     client.Roomid,
			Client: make(map[*Client]bool),
		}
		created_room := ws.ChatRooms[client.Roomid]
		created_room.Client[client] = true
	}
	go ws.Subscribetochannel(ctx, client.Roomid)
}

func (ws *Wserver) unregisterClient(ctx context.Context, client *Client) {
	delete(ws.ChatRooms[client.Roomid].Client, client)

	if len(ws.ChatRooms[client.Roomid].Client) == 0 {
		delete(ws.ChatRooms, client.Roomid)

	} else {

		msg := &Message{
			content: []byte(client.username + "has left the chat"),
			roomid:  client.Roomid,
			sender:  client.username,
		}

		ws.BroadcastMessage(ctx, msg)
	}

}

func (ws *Wserver) BroadcastMessage(ctx context.Context, msg *Message) {

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling the message: %v", err)
	}
	err = ws.Cache.Publish(ctx, msg.roomid.String(), data).Err()
	if err != nil {
		log.Printf("Error publishing message to Redis: %v", err)
	}
}

func (ws *Wserver) Subscribetochannel(ctx context.Context, roomid uuid.UUID) {
	pubsub := ws.Cache.Subscribe(context.Background(), roomid.String())

	defer func() {
		pubsub.Unsubscribe(ctx, roomid.String())
		pubsub.Close()
	}()
	ch := pubsub.Channel()

	for msg := range ch {
		var message Message
		err := json.Unmarshal([]byte(msg.Payload), message)
		if err != nil {
			log.Printf("Error unmarshalling the message", err)
		}
		ws.Publishtoclients(&message)
	}
}

func (ws *Wserver) Publishtoclients(msg *Message) {
	for client := range ws.ChatRooms[msg.roomid].Client {
		if client.username != msg.sender {
			client.Message <- msg
		}
	}
}
