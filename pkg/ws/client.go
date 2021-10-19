package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		fmt.Println(len(message), string(message), err)
		if len(message) > 0 {
			var msg ClientMessage
			errm := json.Unmarshal([]byte(message), &msg)
			if errm != nil {
				log.Println("Can't Unmarshal ClientMessage", errm.Error())
			}
			if msg.Action == "subscribe" {
				c.hub.clients[c].isSubs = true
				c.hub.clients[c].Symbols = msg.Symbols
				fmt.Printf("Client message subscribe: %v \n", c.hub.clients[c])
			} else if msg.Action == "unsubscribe" {
				c.hub.clients[c].isSubs = false
				c.hub.clients[c].Symbols = []string{}
				fmt.Printf("Client message unsubscribe: %v \n", c.hub.clients[c])
			}
		}
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump(bws *WS) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			cSettings, ok := c.hub.clients[c]
			if !ok {
				continue
			}
			fmt.Printf("%s !!!! %v \n", string(message), ok)
			updates := bws.GetChannel()
		Loop:
			for update := range updates {
				if !cSettings.isSubs {
					continue Loop
				}
				var target Update
				if len(cSettings.Symbols) != 0 {
					for _, val := range cSettings.Symbols {
						if val == update.Data[0].Symbol {
							target = update
						} else {
							break Loop
						}
					}
				} else {
					target = update
				}
				if len(target.Data) == 0 {
					continue Loop
				}
				trgObj := UpdateClientMessage{
					Timestamp: target.Data[0].Timestamp,
					Symbol:    target.Data[0].Symbol,
					Price:     target.Data[0].MarkPrice,
				}
				msg, err := json.Marshal(trgObj)
				if err != nil {
					fmt.Printf("Marshal err: %s", err.Error())
				}
				c.conn.WriteMessage(websocket.TextMessage, msg)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
