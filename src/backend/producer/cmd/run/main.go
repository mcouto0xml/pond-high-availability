package main

import (
	"context"
	"fmt"
	"producer/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Erro ao fazer o carregamento do .env: %v", err)
	}
	httpServer := server.NewServer(":8080", &ctx)
	httpServer.Start()
}