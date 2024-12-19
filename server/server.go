package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	Name      string `json:"name"`
	User_name string `json:"user_name"`
	Ip_public string `json:"ip_public"`
	Ip_local  string `json:"ip_local"`
	OS        string `json:"os"`
}

type connClient struct {
	Conn   *websocket.Conn
	Client Client
}

type Clients struct {
	clients map[int]*connClient
	mu      sync.Mutex
}

type Admins struct {
	admins map[int]*connClient
	mu     sync.Mutex
}

var clients Clients

var admins Admins

func addClient(conn *websocket.Conn) (int, error) {
	var client Client
	err := conn.ReadJSON(&client)
	if err != nil {
		return 0, err
	}
	id := 1
	var connClient connClient
	connClient.Client = client
	connClient.Conn = conn
	clients.mu.Lock()
	for {
		if _, ok := clients.clients[id]; !ok {
			clients.clients[id] = &connClient
			break
		}
		id++
	}
	clients.mu.Unlock()
	return id, nil
}

func removeClient(id int) {
	clients.clients[id].Conn.Close()
	clients.mu.Lock()
	delete(clients.clients, id)
	clients.mu.Unlock()

}

func admin(w http.ResponseWriter, r *http.Request) {

}

func client(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error to upgrade client connection:", err)
		return
	}

	id, err := addClient(conn)

	if err != nil {
		conn.Close()
		return
	}

	defer removeClient(id)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	ch1 := make(chan string)
	go func() {
		for {
			select {
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					ch1 <- "error"
					return
				}
			}
		}
	}()

	// ДОДЕЛАТЬ ПОЛУЧЕНИЯ РЕЗУЛЬТАТА ОТ КЛИЕНТА

	select {
	case <-ch1:
		return
	}

}

func main() {
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/client", client)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Panicln("ListenAndServe err:", err)
	}
}
