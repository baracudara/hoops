package ws

import (
	"fmt"
	"sync"
)

type Client struct {
    PlayerID string
    Send     chan []byte
}

type Hub struct {
    clients    map[string]*Client // playerID → client
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]*Client),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
		
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.PlayerID] = client
            h.mu.Unlock()
            fmt.Println("client registered:", client.PlayerID)

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client.PlayerID]; ok {
                delete(h.clients, client.PlayerID)
                close(client.Send)
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            fmt.Println("broadcast received, clients:", len(h.clients))
            h.mu.RLock()
            for _, client := range h.clients {
                select {
                case client.Send <- message:
                default:
                    close(client.Send)
                    delete(h.clients, client.PlayerID)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) SendToPlayer(playerID string, message []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    if client, ok := h.clients[playerID]; ok {
        client.Send <- message
    }
}

func (h *Hub) Broadcast(message []byte) {
    h.broadcast <- message
}