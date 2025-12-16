package main

import (
	"log"

	"goway/pkg/server"
)

func main() {
	srv := server.New()

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
