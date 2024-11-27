package main

import (
	"log"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/server"
)

func main() {
	server := server.NewServer()
	log.Println("Server is running on port", cfg.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
