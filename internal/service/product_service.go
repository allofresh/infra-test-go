package service

import (
	"errors"
	"sync"

	"github.com/allofresh/infra-test-go/internal/model"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
)

// ProductService manages product operations with concurrent-safe access.
type ProductService struct {
	mu       sync.RWMutex
	products map[int]*model.Product
	nextID   int
}

// NewProductService creates a new ProductService instance.
func NewProductService() *ProductService {
	return &ProductService{
		products: make(map[int]*model.Product),
		nextID:   1,
	}
}

// AddProduct adds a new product and returns its assigned ID.
func (s *ProductService) AddProduct(name string, price float64, quantity int) (int, error) {
	if name == "" || price < 0 || quantity < 0 {
		return 0, ErrInvalidProduct
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.products[id] = &model.Product{
		ID:       id,
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}
	s.nextID++
	return id, nil
}

// GetProduct retrieves a product by ID.
func (s *ProductService) GetProduct(id int) (*model.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}
	return p, nil
}

// ListProducts returns all products.
func (s *ProductService) ListProducts() []*model.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*model.Product, 0, len(s.products))
	for _, p := range s.products {
		result = append(result, p)
	}
	return result
}

// UpdateQuantity updates the stock quantity for a product.
func (s *ProductService) UpdateQuantity(id int, quantity int) error {
	if quantity < 0 {
		return ErrInvalidProduct
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[id]
	if !ok {
		return ErrProductNotFound
	}
	p.Quantity = quantity
	return nil
}

// DeleteProduct removes a product by ID.
func (s *ProductService) DeleteProduct(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[id]; !ok {
		return ErrProductNotFound
	}
	delete(s.products, id)
	return nil
}
