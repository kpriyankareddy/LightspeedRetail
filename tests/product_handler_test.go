package tests

import (
	"LightspeedRetail/handlers"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to extract error message from JSON response
func getErrorMessage(resp *httptest.ResponseRecorder) string {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var responseMap map[string]string
	json.Unmarshal(bodyBytes, &responseMap)

	return responseMap["error"] // Extracts "error" field from JSON response
}

// Test GET /products (Success Case)
func TestGetProducts_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "/products", nil)
	resp := httptest.NewRecorder()

	handlers.GetProductsHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

// Test POST /products - Valid Product
func TestCreateProduct_Success(t *testing.T) {
	body := `{"name":"Bluetooth Speaker","price":49.99}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateProductHandler(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.Code)
	}
}

// Test POST /products - Invalid JSON Format
func TestCreateProduct_InvalidJSON(t *testing.T) {
	body := `{invalid-json}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateProductHandler(resp, req)

	expectedError := "Invalid request body"
	actualError := getErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}

// Test POST /products - Missing Name
func TestCreateProduct_MissingName(t *testing.T) {
	body := `{"name":"","price":20.00}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateProductHandler(resp, req)

	expectedError := "Invalid product details: name is required"
	actualError := getErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}

// Test POST /products - Negative Price
func TestCreateProduct_NegativePrice(t *testing.T) {
	body := `{"name":"Faulty Item","price":-20}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateProductHandler(resp, req)

	expectedError := "Invalid product details: price must be greater than zero"
	actualError := getErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}

// Test POST /products - Zero Price
func TestCreateProduct_ZeroPrice(t *testing.T) {
	body := `{"name":"Gift Item","price":0}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateProductHandler(resp, req)

	expectedError := "Invalid product details: price must be greater than zero"
	actualError := getErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}

// Test POST /products - Missing Price
func TestCreateProduct_MissingPrice(t *testing.T) {
	body := `{"name":"Wireless Charger"}` // Price is missing
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handlers.CreateProductHandler(resp, req)

	expectedError := "Invalid product details: price must be greater than zero"
	actualError := getErrorMessage(resp)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
	if actualError != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
	}
}
