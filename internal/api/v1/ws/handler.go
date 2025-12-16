package ws

import (
    "log"

    "github.com/gorilla/websocket"
)

type Handler struct {
    hub *Hub
}

func NewHandler() *Handler {
    h := &Handler{
        hub: newHub(),
    }
    h.hub.start()
    return h
}

func (h *Handler) HandleConnection(conn *websocket.Conn) {
    client := &Client{conn: conn, send: make(chan []byte)}

    h.hub.register <- client

    go h.readPump(client)
    go h.writePump(client)

    client.send <- []byte("Welcome to WebSocket!")
}

type Client struct {
    conn *websocket.Conn
    send chan []byte
}

func (h *Handler) readPump(client *Client) {
    defer func() {
        h.hub.unregister <- client
        close(client.send)
        client.conn.Close()
    }()
    for {
        _, msg, err := client.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Unexpected read error: %v", err)
            }
            break
        }
        h.hub.broadcast <- msg
    }
}

func (h *Handler) writePump(client *Client) {
    defer client.conn.Close()
    for msg := range client.send {
        if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            log.Println("Write error:", err)
            return
        }
    }
}
