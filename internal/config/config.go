package config

import (
    "log"
    "os"
)

type Config struct {
    Port string
    Addr string
}

func Load() *Config {
    cfg := &Config{
        Port: os.Getenv("PORT"),
        Addr: os.Getenv("LISTEN_ADDR"),
    }
    if cfg.Port == "" {
        cfg.Port = "8080"
    }
    if cfg.Addr == "" {
        cfg.Addr = "localhost"
    }
    log.Printf("Config loaded: %+v", cfg)
    return cfg
}
