package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// struct containing message
type Message struct { //struct with message and recipient
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	endpoint := "ws://localhost:8080/ws?username=" + username   //server address
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil) //connecting to server
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close() //closing connection
	go func() {        //goroutine function
		for { //infinite loop
			var msg Message            //message struct
			err := conn.ReadJSON(&msg) //converts the json to text
			if err != nil {
				log.Fatalf("Failed to read message from WebSocket server: %v", err)
			}
			fmt.Printf("\nReceived message from %s: %s\n", msg.Sender, msg.Text)
		}
	}()
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin) //standard input
			fmt.Print("Enter recipient's username: ")
			recipient, _ := reader.ReadString('\n')
			recipient = strings.TrimSpace(recipient)
			fmt.Print("Enter message: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)
			msg := Message{Sender: username, Recipient: recipient, Text: message} //converting to msg struct
			err := conn.WriteJSON(msg)                                            //converting struct to json
			if err != nil {
				log.Fatalf("Failed to send message %v", err)
			}
		}
	}()
	select {}
}
