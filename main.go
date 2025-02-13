package main

import (
	"LightspeedRetail/handlers"
	"log"
	"net/http"
)

func productRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlers.CreateProductHandler(w, r)
		return
	}
	handlers.GetProductsHandler(w, r)
}

func main() {
	serverAddr := "127.0.0.1:8080"
	http.HandleFunc("/products", productRouter)

	log.Printf("Server started on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
