package services

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

type Message struct {
	Id      string
	Message string `json:"message"`
}

// For broadcast type socket; keep as an example
// type Hub struct {
// 	clients   map[*websocket.Conn]bool
// 	broadcast chan Message
// }

// func NewHub() *Hub {
// 	return &Hub{
// 		clients:   make(map[*websocket.Conn]bool),
// 		broadcast: make(chan Message),
// 	}
// }

// func (h *Hub) run() {
// 	for {
// 		select {
// 		case message := <-h.broadcast:
// 			for client := range h.clients {
// 				fmt.Println(client)
// 				if err := client.WriteJSON(message); !errors.Is(err, nil) {
// 					log.Printf("error occurred: %v", err)
// 				}
// 			}
// 		}
// 	}
// }

// Run will execute Go Routines to check incoming Socket events
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			HandleUserRegisterEvent(hub, client)

		case client := <-hub.unregister:
			HandleUserDisconnectEvent(hub, client)
		}
	}
}
