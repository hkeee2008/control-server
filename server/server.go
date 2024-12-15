package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func admin(w http.ResponseWriter, r *http.Request) {

}
