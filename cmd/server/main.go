package main

import (
	"log"

	"steamshopemulator/internal/app"
)

func main() {
	server, err := app.New()
	if err != nil {
		log.Fatalf("bootstrap server: %v", err)
	}

	log.Printf("steam shop emulator listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
