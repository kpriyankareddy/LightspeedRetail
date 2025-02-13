package tests

import (
	"LightspeedRetail/handlers"
	"LightspeedRetail/repositories"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to extract error message
func getSalesErrorMessage(resp *httptest.ResponseRecorder) string {
	var responseMap map[string]string
	json.Unmarshal(resp.Body.Bytes(), &responseMap)
	return responseMap["error"]
}

// Test POST /sales - Valid Sale
func TestCreateSale_Success(t *testing.T) {
	// Get a real product ID from the repository
	products := repositories.GetAllProducts()
	if len(products) == 0 {
		t.Fatal("No products found in repository")
	}
	productID := products[0].ID // Use the first product

	// Build JSON with a valid product ID
	body := `{"items":[{"product_id":"` + productID.String() + `","quantity":2}]}`

	req, _ := http.NewRequest("POST", "/sales", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateSaleHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

// Test POST /sales - Successful Sale with 2 Items and Different Quantities
func TestCreateSale_SuccessWithMultipleItems(t *testing.T) {
	products := repositories.GetAllProducts()
	if len(products) < 2 {
		t.Fatal("Not enough products in repository for this test")
	}

	product1 := products[0]
	product2 := products[1]

	// Build JSON with multiple items
	body := `{"items":[
		{"product_id":"` + product1.ID.String() + `","quantity":2},
		{"product_id":"` + product2.ID.String() + `","quantity":3}
	]}`

	req, _ := http.NewRequest("POST", "/sales", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateSaleHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Parse response
	var saleResponse map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &saleResponse)

	// Validate total price
	expectedTotal := (2 * product1.Price) + (3 * product2.Price)
	actualTotal := saleResponse["total"].(float64)

	if actualTotal != expectedTotal {
		t.Errorf("Expected total price %.2f, got %.2f", expectedTotal, actualTotal)
	}
}

// Test POST /sales - No Items in Sale
func TestCreateSale_NoItems(t *testing.T) {
	body := `{"items":[]}`
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateSaleHandler(resp, req)

	expectedError := "Sale must contain at least one item"
	actualError := getSalesErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}

// Test POST /sales - Invalid Product ID
func TestCreateSale_InvalidProductID(t *testing.T) {
	body := `{"items":[{"product_id":"00000000-0000-0000-0000-000000000000","quantity":1}]}`
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateSaleHandler(resp, req)

	expectedError := "Product not found"
	actualError := getSalesErrorMessage(resp)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}

// Test POST /sales - Negative Quantity
func TestCreateSale_NegativeQuantity(t *testing.T) {
	products := repositories.GetAllProducts()
	if len(products) == 0 {
		t.Fatal("No products found in repository")
	}
	productID := products[0].ID

	body := `{"items":[{"product_id":"` + productID.String() + `","quantity":-1}]}`
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateSaleHandler(resp, req)

	expectedError := "Quantity must be a positive integer"
	actualError := getSalesErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}

// Test POST /sales - Decimal Quantity
func TestCreateSale_DecimalQuantity(t *testing.T) {
	products := repositories.GetAllProducts()
	if len(products) == 0 {
		t.Fatal("No products found in repository")
	}
	productID := products[0].ID

	body := `{"items":[{"product_id":"` + productID.String() + `","quantity":1.5}]}`
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateSaleHandler(resp, req)

	expectedError := "Invalid request body" // Since Go rejects float in int fields
	actualError := getSalesErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}
