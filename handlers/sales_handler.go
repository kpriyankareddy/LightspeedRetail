package handlers

import (
	"LightspeedRetail/models"
	"LightspeedRetail/repositories"
	"encoding/json"
	"math"
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

	totalPrice := 0.0

	// Process each line item
	for i := range sale.Items {
		// Validate that quantity is a positive integer FIRST
		if sale.Items[i].Quantity <= 0 {
			sendSalesErrorResponse(w, http.StatusBadRequest, "Quantity must be a positive integer")
			return
		}

		product, err := repositories.GetProductByID(sale.Items[i].ProductID)
		if err != nil {
			sendSalesErrorResponse(w, http.StatusNotFound, "Product not found")
			return
		}

		// Calculate total for this item
		sale.Items[i].Total = float64(sale.Items[i].Quantity) * product.Price
		totalPrice += sale.Items[i].Total
	}

	// Validate discount AFTER quantity check
	if sale.Discount < 0 {
		sendSalesErrorResponse(w, http.StatusBadRequest, "Discount must be a positive value")
		return
	}

	if sale.Discount > totalPrice {
		sendSalesErrorResponse(w, http.StatusBadRequest, "Discount cannot exceed total price")
		return
	}

	// Distribute discount proportionally
	remainingDiscount := sale.Discount
	for i := range sale.Items {
		itemShare := (sale.Items[i].Total / totalPrice) * sale.Discount
		sale.Items[i].Discount = math.Floor(itemShare*100) / 100 // Round down to 2 decimal places
		remainingDiscount -= sale.Items[i].Discount
	}

	// Adjust last item to fix rounding errors
	lastIndex := len(sale.Items) - 1
	sale.Items[lastIndex].Discount += remainingDiscount

	// Update the total after applying the discount
	sale.Total = totalPrice - sale.Discount

	// Return the updated sale response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}
