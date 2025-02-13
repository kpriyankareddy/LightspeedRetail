package main

import (
	"LightspeedRetail/handlers"
	"log"
	"net/http"
)

func main() {
	serverAddr := "127.0.0.1:8080"
	http.HandleFunc("/products", handlers.GetProductsHandler)

	log.Printf("Server started on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
