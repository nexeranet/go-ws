package ws

type ClientSettings struct {
	isRegistred bool
	Symbols     []string
	isSubs      bool
}

type Hub struct {
	// Registered clients.
	clients map[*Client]*ClientSettings

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]*ClientSettings),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			cl, ok := h.clients[client]
			if ok {
				cl.isRegistred = true
			} else {
				h.clients[client] = &ClientSettings{
					isRegistred: true,
					Symbols:     []string{},
					isSubs:      false,
				}
			}
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
