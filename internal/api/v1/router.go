package v1

import (
	"time"
	"net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "goway/internal/api/v1/ws"
)

func Router() http.Handler {
    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(30 * time.Second))

    r.Mount("/ws", ws.Routes())

    return r
}
