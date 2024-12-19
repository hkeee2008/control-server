package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

var addr = flag.String(os.Args[1], "localhost:"+os.Args[2], "http service address")

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/client"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Panicln("Error connection to server:", err)
	}
	fmt.Println(u.String())
	defer conn.Close()

	_, msg, err := conn.ReadMessage()

	if err != nil {
		log.Panicln("Error to read message from server:", err)
	}

	fmt.Println(string(msg))
}
