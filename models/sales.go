package models

import "github.com/google/uuid"

// SaleLineItem represents a single item in a sale
type SaleLineItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Discount  float64   `json:"discount"`
}

// Sale represents a sales request
type Sale struct {
	Items    []SaleLineItem `json:"items"`
	Total    float64        `json:"total"`
	Discount float64        `json:"discount"`
}
