package main

import (
	"fmt"
	"net/http"
	"consumer/internal/function"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /telemetry", function.PostTelemetry)

	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
	fmt.Print("Escutando na porta :8080")
}