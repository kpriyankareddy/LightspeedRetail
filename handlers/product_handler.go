package handlers

import (
	"LightspeedRetail/models"
	"LightspeedRetail/repositories"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// Helper function to send JSON error responses
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// Get all the products
func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	products := repositories.GetAllProducts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Create a new product
func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	// Read and decode JSON request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		sendErrorResponse(w, http.StatusBadRequest, "Error reading request body")
		return
	}
	defer r.Body.Close()

	log.Println("Received Request Body:", string(body))

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		log.Println("Invalid JSON:", err)
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if product.Name == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid product details: name is required")
		return
	}
	if product.Price <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid product details: price must be greater than zero")
		return
	}

	// Assign a unique ID
	product.ID = uuid.New()
	repositories.AddProduct(product)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
