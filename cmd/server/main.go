package main

import (
    "context"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"

    "goway/internal/api/v1"
    "goway/internal/config"
)

func main() {
    cfg := config.Load()

    listenAddr := net.JoinHostPort(cfg.Addr, cfg.Port)

    r := chi.NewRouter()

    r.Get("/test-ws", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "testdata/index.html")
    })

    r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })

    r.Mount("/api/v1", v1.Router())

    srv := &http.Server{
        Addr:    listenAddr,
        Handler: r,
    }

    go func() {
        log.Printf("Server starting on %s", listenAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe error: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutdown signal received")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown error: %v", err)
    }
    log.Println("Server gracefully stopped")
}
