package ws

import (
    "log"

    "github.com/gorilla/websocket"
)

type Handler struct {
}

func NewHandler() *Handler {
    return &Handler{}
}

func (h *Handler) HandleConnection(conn *websocket.Conn) {
    client := &Client{conn: conn, send: make(chan []byte)}

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
        close(client.send)
        client.conn.Close()
    }()
    for {
        _, msg, err := client.conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            break
        }
        client.send <- msg
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
