package tests

import (
	"LightspeedRetail/handlers"
	"LightspeedRetail/repositories"
	"bytes"
	"encoding/json"
	"math"
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

func TestCreateSale_WithValidDiscount(t *testing.T) {
	// Get actual products from repository to ensure valid product IDs
	products := repositories.GetAllProducts()
	if len(products) < 2 {
		t.Fatal("Not enough products in repository for this test")
	}

	product1 := products[0]
	product2 := products[1]

	// Create a sale with two items and a discount
	body := `{
		"items": [
			{"product_id": "` + product1.ID.String() + `", "quantity": 2},
			{"product_id": "` + product2.ID.String() + `", "quantity": 3}
		],
		"discount": 20.00
	}`

	req, _ := http.NewRequest("POST", "/sales", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Call the sale handler
	handlers.CreateSaleHandler(resp, req)

	// Ensure the response is successful
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Decode response safely
	var response map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Ensure "total" and "discount" exist in the response before using them
	if _, exists := response["total"]; !exists {
		t.Fatalf("Response does not contain 'total' field: %v", response)
	}
	if _, exists := response["discount"]; !exists {
		t.Fatalf("Response does not contain 'discount' field: %v", response)
	}

	// Expected values
	expectedTotal := (float64(2)*product1.Price + float64(3)*product2.Price) - 20.00
	actualTotal := response["total"].(float64)

	if actualTotal != expectedTotal {
		t.Errorf("Expected total after discount %.2f, got %.2f", expectedTotal, actualTotal)
	}

	// Check if the discount is distributed correctly
	items := response["items"].([]interface{})
	item1 := items[0].(map[string]interface{})
	item2 := items[1].(map[string]interface{})

	expectedItem1Discount := (float64(2) * product1.Price / (float64(2)*product1.Price + float64(3)*product2.Price)) * 20.00
	expectedItem2Discount := (float64(3) * product2.Price / (float64(2)*product1.Price + float64(3)*product2.Price)) * 20.00

	actualItem1Discount := item1["discount"].(float64)
	actualItem2Discount := item2["discount"].(float64)

	// Allow minor rounding differences due to float precision
	const tolerance = 0.01
	if math.Abs(actualItem1Discount-expectedItem1Discount) > tolerance {
		t.Errorf("Expected item1 discount %.2f, got %.2f", expectedItem1Discount, actualItem1Discount)
	}
	if math.Abs(actualItem2Discount-expectedItem2Discount) > tolerance {
		t.Errorf("Expected item2 discount %.2f, got %.2f", expectedItem2Discount, actualItem2Discount)
	}
}
