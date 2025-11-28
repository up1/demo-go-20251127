package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []byte)
	rdb       = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx = context.Background()
)

// Handle WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Publish message to Redis Channel
		if err := rdb.Publish(ctx, "websocketChannel", msg).Err(); err != nil {
			log.Println("Error publishing to Redis:", err)
		}
	}
}

// Handle messages from Redis and broadcast to all clients
func handleRedisMessages() {
	sub := rdb.Subscribe(ctx, "websocketChannel")
	ch := sub.Channel()

	for msg := range ch {
		log.Printf("Received message from Redis: %s", msg.Payload)
		broadcast <- []byte(msg.Payload)
	}
}

// Broadcast messages to all clients
func handleBroadcasts() {
	for {
		msg := <-broadcast
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("Error writing message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Check redis connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully")

	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/", serveHome)

	go handleBroadcasts()
	go handleRedisMessages()

	log.Println("WebSocket server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}
