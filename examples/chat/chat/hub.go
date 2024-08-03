package chat

var hub *Hub

func init() {
	hub = NewHub()
	go hub.Run()
}

type Hub struct {
	clients    map[*Client]struct{} // registered clients
	broadcast  chan []byte          // inbound messages from the clients
	register   chan *Client         // register requests from the clients
	unregister chan *Client         // unregister requests from clients
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
