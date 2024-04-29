package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients = make(map[*websocket.Conn]*Client) //map of all clients
	mu      sync.Mutex                          //ensures that only one goroutine can access the clients map
)

type Client struct { //struct with connection and username
	conn     *websocket.Conn
	username string
}
type Message struct { //struct with message and recipient
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

func websockethandler(w http.ResponseWriter, r *http.Request) { //websocket handler
	conn, err := upgrader.Upgrade(w, r, nil) //connection
	if err != nil {
		log.Println("Failed to upgrade connection to WebSocket:", err)
		return
	}
	defer conn.Close()                        //closes the connection after use
	username := r.URL.Query().Get("username") //looks for username in the request
	if username == "" {
		log.Println("Username not provided")
		return
	}
	client := &Client{
		conn:     conn,
		username: username,
	} //makes new client
	mu.Lock()              //allows change in clients map
	clients[conn] = client //adds user connection to clients map
	mu.Unlock()            //closes map
	for {
		var msg Message
		err := conn.ReadJSON(&msg) //reads the json sent by client
		if err != nil {
			log.Println("Error reading msg from WebSocket:", err)
			break
		}
		mu.Lock()
		recipientConn := findConn(msg.Recipient) //adds the recipient of the message to clients map
		mu.Unlock()
		if recipientConn != nil {
			log.Printf("Recipient '%s' not found", msg.Recipient)
			continue
		}
		err = recipientConn.WriteJSON(msg) //sends message to the recipient in form of json
		if err != nil {
			log.Println("Error sending msg to recipient:", err)
			continue
		}
	}
}
func findConn(username string) *websocket.Conn {
	for conn, client := range clients {
		if client.username == username {
			return conn
		}
	}
	return nil
}

func main() {
	http.HandleFunc("/ws", websockethandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
