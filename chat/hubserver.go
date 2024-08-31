package main

func (ws *Wserver) Runserver() {

	for {
		select {
		case Client := <-ws.Register:
			ws.registerClient(Client)
		case Client := <-ws.Unregister:
			ws.unregisterClient(Client)
		case Message := <-ws.Broadcast:
			ws.BroadcastMessage(Message)
		}

	}
}

func (ws *Wserver) registerClient(client *Client) {
	room, ok := ws.ChatRooms[client.Roomid]
	if ok {
		room.Client[client] = true
		msg := &Message{
			content: []byte(client.username + "has joined the chat"),
			roomid:  client.Roomid,
			sender:  client.username,
		}
		client.Message <- msg
	}
	if !ok {
		ws.ChatRooms[client.Roomid] = &chatRooms{
			ID:     client.Roomid,
			Client: make(map[*Client]bool),
		}
		created_room := ws.ChatRooms[client.Roomid]
		created_room.Client[client] = true
	}
}

func (ws *Wserver) unregisterClient(client *Client) {
	if len(ws.ChatRooms[client.Roomid].Client) == 0 {
		delete(ws.ChatRooms[client.Roomid].Client, client)
		delete(ws.ChatRooms, client.Roomid)
	}

	msg := &Message{
		content: []byte(client.username + "has left the chat"),
		roomid:  client.Roomid,
		sender:  client.username,
	}
	client.ws.Broadcast <- msg
	delete(ws.ChatRooms[client.Roomid].Client, client)

}

func (ws *Wserver) BroadcastMessage(msg *Message) {
	for client := range ws.ChatRooms[msg.roomid].Client {
		if client.username != msg.sender {
			client.Message <- msg
		}
	}
}
