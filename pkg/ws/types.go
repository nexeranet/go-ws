package ws

import "github.com/gorilla/websocket"

type UpdateChannel chan Update

type WS struct {
	message_channel  UpdateChannel
	shutdown_channel chan interface{}
	ws               *websocket.Conn
}

func NewWS() *WS {
	return &WS{
		shutdown_channel: make(chan interface{}),
	}
}

type AuthMessage struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type CommandMessage struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

type Data struct {
	Symbol    string  `json:"symbol"`
	MarkPrice float64 `json:"markPrice"`
}

type Update struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   []Data `json:"data"`
}
