package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocket server struct
type Server struct {
	clients          map[string]*websocket.Conn
	clientsLock      sync.Mutex
	upgrader         websocket.Upgrader
}


func (s *Server) handleSpell(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Query().Get("id"))
}

// Handle WebSocket connections
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
 	
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Generate a unique URL for the client
	url := uuid.New().String()

	// Store the connection URL and WebSocket object
	s.clientsLock.Lock()
	s.clients[url] = conn
	s.clientsLock.Unlock()

	// Send the URL to the client
	if err := conn.WriteMessage(websocket.TextMessage, []byte(url)); err != nil {
		log.Println("Error sending URL to client:", err)
		return
	}

	// Handle incoming messages from the client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client:", err)
			break
		}
		log.Println("Message received:", string(message))
		// Broadcast the message to all subscribed clients
		s.broadcastMessage(url, message)
	}

	// Remove the connection URL and WebSocket object
	s.clientsLock.Lock()
	delete(s.clients, url)	
	s.clientsLock.Unlock()
}

// Broadcast message to all connected clients
func (s *Server) broadcastMessage(senderURL string, message []byte) {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	for url, conn := range s.clients {
		if url != senderURL {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error broadcasting message:", err)
				continue
			}
		}
	}
}

func main() {
	server := &Server{
		clients: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	http.HandleFunc("/spell", server.handleSpell)
	http.HandleFunc("/", server.handleConnection)

	fmt.Println("WebSocket server started on port 8765")
	log.Fatal(http.ListenAndServe(":8765", nil))
}
