package repositories

import (
	"LightspeedRetail/models"
	"errors"
	"sync"

	"github.com/google/uuid"
)

// In-memory product storage
var (
	mu       sync.Mutex
	products = []models.Product{
		{ID: uuid.New(), Name: "Chrome Toaster", Price: 100},
		{ID: uuid.New(), Name: "Copper Kettle", Price: 49.99},
		{ID: uuid.New(), Name: "Mixing Bowl", Price: 20},
	}
)

// Returns all stored products
func GetAllProducts() []models.Product {
	mu.Lock()
	defer mu.Unlock()
	return append([]models.Product{}, products...)
}

// Adds a new product to the list
func AddProduct(product models.Product) {
	mu.Lock()
	defer mu.Unlock()
	products = append(products, product)
}

// Retrieves a product by ID
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
