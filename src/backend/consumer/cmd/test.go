// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"function.com/consumer/function"
// )

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	http.HandleFunc("/", function.PostTelemetry)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }

package test