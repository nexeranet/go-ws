package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	go_ws "github.com/nexeranet/go-ws"
	"github.com/nexeranet/go-ws/pkg/handler"
	"github.com/nexeranet/go-ws/pkg/ws"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// start
	WS := ws.NewWS()
	err = WS.Start()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	defer WS.Close()

	handler := handler.NewHandler()
	fmt.Println(handler)
	srv := new(go_ws.Server)
	fmt.Printf("%v", srv)
	err = srv.Run("8080", handler.InitRouter(WS))
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
