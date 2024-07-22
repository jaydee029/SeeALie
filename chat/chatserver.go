package main

type wsserver struct {
	Clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func newWebServer() *wsserver {
	return &wsserver{
		Clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (ws *wsserver) registerClient(client *Client) {
	ws.Clients[client] = true
}

func (ws *wsserver) unregisterClient(client *Client) {
	delete(ws.Clients, client)
}

func (ws *wsserver) Run() {

	for {
		select {
		case Client := <-ws.register:
			ws.registerClient(Client)

		case Client := <-ws.unregister:
			ws.unregisterClient(Client)
		}
	}
}
