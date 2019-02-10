package main

import (
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello gophercon!")
}

func ws(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// Read messages from socket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message %v", err)
			return
		}
		log.Println(string(msg))
	}
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/ws", ws)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
