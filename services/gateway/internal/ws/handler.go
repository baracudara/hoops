package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"

	custommiddleware "github.com/baracudara/hoops/gateway/internal/middleware"
	"github.com/baracudara/hoops/protos/gen/go/auth"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // для разработки разрешаем все origins
    },
}

type Handler struct {
    hub *Hub
    log *slog.Logger
	ticketStore *TicketStore
}

func NewHandler(hub *Hub, log *slog.Logger, ticketStore *TicketStore) *Handler {
    return &Handler{
		hub: hub, 
		log: log, 
		ticketStore: ticketStore}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
    ticket := r.URL.Query().Get("ticket")
    if ticket == "" {
        http.Error(w, "ticket required", http.StatusUnauthorized)
        return
    }

    playerID, err := h.ticketStore.Consume(r.Context(), ticket)
    if err != nil {
        http.Error(w, "invalid ticket", http.StatusUnauthorized)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        h.log.Error("failed to upgrade connection", "err", err)
        return
    }

    client := &Client{
        PlayerID: playerID,
        Send:     make(chan []byte, 256),
    }

    h.hub.register <- client

    go client.writePump(conn)
    go client.readPump(conn, h.hub)
}


func (c *Client) writePump(conn *websocket.Conn) {
    defer conn.Close()

    for {
        message, ok := <-c.Send
        if !ok {
            // канал закрыт — отправляем close и выходим
            conn.WriteMessage(websocket.CloseMessage, []byte{})
            return
        }

        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            return
        }
    }
}

func (c *Client) readPump(conn *websocket.Conn, hub *Hub) {
    defer func() {
        hub.unregister <- c
        conn.Close()
    }()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            // клиент отключился
            break
        }
    }
}

func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
    playerID := r.Context().Value(custommiddleware.UserKey).(*auth.VerifyAccessTokenResponse)

    ticket, err := h.ticketStore.Create(r.Context(), playerID.Uuid)
    if err != nil {
        h.log.Error("failed to create ticket", "err", err)
        http.Error(w, "failed to create ticket", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "ticket": ticket,
    })
}