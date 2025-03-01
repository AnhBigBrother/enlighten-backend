package main

import (
	"fmt"
	"log"
	"net"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/server"
)

func main() {
	// server := server.NewServer()
	// log.Println("Server is running on port", cfg.Port)
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	server := server.NewGrpcServer()
	address := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	log.Println("Start server on port", cfg.Port)

	err = server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
