package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Neil-uli/go-play/channels/chat/models"
	"github.com/gorilla/websocket"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true

	for {
		var msg models.ChatMessage
		// Read in a new message as JSON and map it to a Message object
		// is checking to see if msg is populated. If msg is ever not nil, it'll send the message over to the broadcaster channel
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		// send new message to the channel
		broadcaster <- msg
	}
}

func sendPreviousMessages(ws *websocket.Conn) {
	chatMessages, err := rdb.LRange("chat_messages", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	// send previous messages
	for _, chatMessage := range chatMessages {
		var msg models.ChatMessage
		json.Unmarshal([]byte(chatMessage), &msg)
		err := client.WriteJSON(msg)
		if err != nil && unsafeError(err) {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
func handleMessages() {
	for {
		// grab any next message from channel
		msg := <-broadcaster
	}
}

func messageClients(msg models.ChatMessage) {
	// send to every client currently connected
	for client := range clients {
		messageClient(client, msg)
	}
}
func messageClient(client *websocket.Conn, msg models.ChatMessage) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error : %v", err)
		client.Close()
		delete(clients, client)
	}
}