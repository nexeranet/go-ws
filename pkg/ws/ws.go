package ws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	VERB     = "GET"
	ENDPOINT = "/realtime"
)

func WsHandler(w http.ResponseWriter, r *http.Request, pWS *WS, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump(pWS)
	go client.readPump()
}

func createSignature(api_secret, verb, endpoint string, exp int64) string {
	message := (verb + endpoint + fmt.Sprintf("%d", exp))
	fmt.Printf("Signing: %s", message)
	h := hmac.New(sha256.New, []byte(api_secret))
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func (obj *WS) Send(message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = obj.ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (obj *WS) Auth() error {
	expires := time.Now().Unix() + 5
	signature := createSignature(os.Getenv("API_SECRET"), VERB, ENDPOINT, expires)
	authMessage := AuthMessage{
		Op:   "authKeyExpires",
		Args: []interface{}{os.Getenv("API_KEY"), expires, signature},
	}
	err := obj.Send(authMessage)
	if err != nil {
		return err
	}
	return nil
}

func (obj *WS) Close() error {
	return obj.ws.Close()
}
func (obj *WS) CloseMessage() error {
	err := obj.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return err
	}
	return nil
}

func (obj *WS) GetChannel() UpdateChannel {
	return obj.message_channel
}

func (obj *WS) StopUpdateChannel() {
	close(obj.shutdown_channel)
}

func (obj *WS) Start() error {
	obj.message_channel = make(UpdateChannel)
	url := os.Getenv("BITMEX_URL")
	log.Printf("connecting to %s", url)
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	obj.ws = ws
	go func() {
		select {
		case <-obj.shutdown_channel:
			close(obj.message_channel)
			return
		default:
		}
		for {
			_, message, err := obj.ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var msg Update
			//log.Printf("Message: %s, $$$$$$$$$$$$$$$$$$$$$$$$$$\n", message)
			err = json.Unmarshal([]byte(message), &msg)
			if err != nil {
				log.Println("Can't Unmarshal", err)
				return
			}
			obj.message_channel <- msg
		}
	}()
	err = obj.Auth()
	if err != nil {
		return err
	}
	msg := CommandMessage{
		Op:   "subscribe",
		Args: []string{"instrument"},
	}
	if err = obj.Send(msg); err != nil {
		return err
	}
	return nil
}
