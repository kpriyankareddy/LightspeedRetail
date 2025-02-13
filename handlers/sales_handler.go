package handlers

import (
	"LightspeedRetail/models"
	"LightspeedRetail/repositories"
	"encoding/json"
	"net/http"
)

// Error response structure for sales
type SalesErrorResponse struct {
	Error string `json:"error"`
}

// Helper function to send JSON error responses
func sendSalesErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SalesErrorResponse{Error: message})
}

// CreateSaleHandler processes a sale transaction
func CreateSaleHandler(w http.ResponseWriter, r *http.Request) {
	var sale models.Sale

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&sale); err != nil {
		sendSalesErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate sale items
	if len(sale.Items) == 0 {
		sendSalesErrorResponse(w, http.StatusBadRequest, "Sale must contain at least one item")
		return
	}

	total := 0.0

	// Process each line item
	for i, item := range sale.Items {
		// Validate quantity (must be positive integer)
		if item.Quantity <= 0 {
			sendSalesErrorResponse(w, http.StatusBadRequest, "Quantity must be a positive integer")
			return
		}

		// Fetch product details
		product, err := repositories.GetProductByID(item.ProductID)
		if err != nil {
			sendSalesErrorResponse(w, http.StatusNotFound, "Product not found")
			return
		}

		// Calculate total for this item
		item.Total = float64(item.Quantity) * product.Price
		sale.Items[i] = item
		total += item.Total
	}

	// Set total sale price
	sale.Total = total

	// Return the sale response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}
