package tests

import (
	"LightspeedRetail/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProducts_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "/products", nil)
	resp := httptest.NewRecorder()

	handlers.GetProductsHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}
