package ws

import "log"

type Hub struct {
    clients     map[*Client]bool
    broadcast   chan struct {
        room string
        msg  []byte
    }
    rooms       map[string]map[*Client]bool
    register    chan *Client
    unregister  chan *Client
    subscribe   chan struct {
        client *Client
        room   string
    }
    unsubscribe chan struct {
        client *Client
        room   string
    }
}

func newHub() *Hub {
    return &Hub{
        clients:     make(map[*Client]bool),
        broadcast:   make(chan struct {
            room string
            msg  []byte
        }, 256),
        rooms:       make(map[string]map[*Client]bool),
        register:    make(chan *Client),
        unregister:  make(chan *Client),
        subscribe:   make(chan struct {
            client *Client
            room   string
        }),
        unsubscribe: make(chan struct {
            client *Client
            room   string
        }),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Println("Registered global client")
        case client := <-h.unregister:
            delete(h.clients, client)
            for room := range h.rooms {
                delete(h.rooms[room], client)
            }
            log.Println("Unregistered client")
        case sub := <-h.subscribe:
            if h.rooms[sub.room] == nil {
                h.rooms[sub.room] = make(map[*Client]bool)
            }
            h.rooms[sub.room][sub.client] = true
            log.Printf("Subscribed client to room %s (total in room: %d)", sub.room, len(h.rooms[sub.room]))
            h.broadcast <- struct {
                room string
                msg  []byte
            }{room: sub.room, msg: []byte("A user joined the room: " + sub.room)}
        case unsub := <-h.unsubscribe:
            delete(h.rooms[unsub.room], unsub.client)
            if len(h.rooms[unsub.room]) == 0 {
                delete(h.rooms, unsub.room)
            }
            log.Printf("Unsubscribed client from room %s (remaining: %d)", unsub.room, len(h.rooms[unsub.room]))
            h.broadcast <- struct {
                room string
                msg  []byte
            }{room: unsub.room, msg: []byte("A user left the room: " + unsub.room)}
        case message := <-h.broadcast:
            roomClients := h.rooms[message.room]
            log.Printf("Broadcasting to room %s (clients: %d, msg: %s)", message.room, len(roomClients), string(message.msg))
            for client := range roomClients {
                select {
                case client.send <- message.msg:
                    log.Println("Sent to client")
                default:
                    log.Println("Dropped slow client from room " + message.room)
                    delete(roomClients, client)
                    delete(h.clients, client)
                }
            }
        }
    }
}

func (h *Hub) start() {
    go h.run()
}

func (h *Hub) stop() {
    close(h.broadcast)
    close(h.register)
    close(h.unregister)
    close(h.subscribe)
    close(h.unsubscribe)
}
