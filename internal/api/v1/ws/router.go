package ws

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/gorilla/websocket"
)

var globalHub = newHub()

func init() {
    globalHub.start()
}

func Routes() http.Handler {
    r := chi.NewRouter()

    var upgrader = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("Upgrade error:", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        handler := NewHandler(globalHub)
        handler.HandleConnection(conn)
    })

    return r
}
