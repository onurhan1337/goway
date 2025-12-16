package ws

import (
    "encoding/json"
    "log"

    "github.com/gorilla/websocket"
)

type Handler struct {
    hub *Hub
}

func NewHandler(hub *Hub) *Handler {
    return &Handler{
        hub: hub,
    }
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

        var message Message
        if err := json.Unmarshal(msg, &message); err != nil {
            log.Printf("An error occurred while parsing JSON: %v", err)
            continue
        }

        switch message.Action {
        case "subscribe":
            h.hub.subscribe <- struct {
                client *Client
                room   string
            }{client: client, room: message.Room}
        case "unsubscribe":
            h.hub.unsubscribe <- struct {
                client *Client
                room   string
            }{client: client, room: message.Room}
        case "send":
            if message.Room != "" && message.Content != "" {
                h.hub.broadcast <- struct {
                    room string
                    msg  []byte
                }{room: message.Room, msg: []byte(message.Content)}
            }
        default:
            log.Println("Unknown action:", message.Action)
        }
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
