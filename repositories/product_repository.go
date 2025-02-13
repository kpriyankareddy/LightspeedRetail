package repositories

import (
	"errors"
	"LightspeedRetail/models"
	"github.com/google/uuid"
	"sync"
)

var mu sync.Mutex
var products = []models.Product{
	{ID: uuid.New(), Name: "Chrome Toaster", Price: 100},
	{ID: uuid.New(), Name: "Copper Kettle", Price: 49.99},
	{ID: uuid.New(), Name: "Mixing Bowl", Price: 20},
}

func GetAllProducts() []models.Product {
	mu.Lock()
	defer mu.Unlock()
	return products
}

func GetProductByID(id uuid.UUID) (*models.Product, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, product := range products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, errors.New("product not found")
}
