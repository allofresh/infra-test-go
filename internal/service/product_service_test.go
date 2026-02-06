package service

import (
	"sync"
	"testing"
)

func TestProductService_AddProduct(t *testing.T) {
	svc := NewProductService()

	id, err := svc.AddProduct("Apple", 1.50, 100)
	if err != nil {
		t.Fatalf("AddProduct() unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("AddProduct() id = %d, want 1", id)
	}

	id2, err := svc.AddProduct("Banana", 2.00, 50)
	if err != nil {
		t.Fatalf("AddProduct() unexpected error: %v", err)
	}
	if id2 != 2 {
		t.Errorf("AddProduct() id = %d, want 2", id2)
	}
}

func TestProductService_AddProduct_Validation(t *testing.T) {
	svc := NewProductService()

	tests := []struct {
		name     string
		prodName string
		price    float64
		quantity int
	}{
		{"empty name", "", 1.00, 10},
		{"negative price", "Test", -1.00, 10},
		{"negative quantity", "Test", 1.00, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.AddProduct(tt.prodName, tt.price, tt.quantity)
			if err != ErrInvalidProduct {
				t.Errorf("AddProduct() error = %v, want %v", err, ErrInvalidProduct)
			}
		})
	}
}

func TestProductService_GetProduct(t *testing.T) {
	svc := NewProductService()
	svc.AddProduct("Apple", 1.50, 100)

	p, err := svc.GetProduct(1)
	if err != nil {
		t.Fatalf("GetProduct() unexpected error: %v", err)
	}
	if p.Name != "Apple" {
		t.Errorf("GetProduct() name = %s, want Apple", p.Name)
	}
	if p.Price != 1.50 {
		t.Errorf("GetProduct() price = %f, want 1.50", p.Price)
	}
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	svc := NewProductService()

	_, err := svc.GetProduct(999)
	if err != ErrProductNotFound {
		t.Errorf("GetProduct() error = %v, want %v", err, ErrProductNotFound)
	}
}

func TestProductService_ListProducts(t *testing.T) {
	svc := NewProductService()

	products := svc.ListProducts()
	if len(products) != 0 {
		t.Errorf("ListProducts() length = %d, want 0", len(products))
	}

	svc.AddProduct("Apple", 1.50, 100)
	svc.AddProduct("Banana", 2.00, 50)

	products = svc.ListProducts()
	if len(products) != 2 {
		t.Errorf("ListProducts() length = %d, want 2", len(products))
	}
}

func TestProductService_UpdateQuantity(t *testing.T) {
	svc := NewProductService()
	svc.AddProduct("Apple", 1.50, 100)

	err := svc.UpdateQuantity(1, 200)
	if err != nil {
		t.Fatalf("UpdateQuantity() unexpected error: %v", err)
	}

	p, _ := svc.GetProduct(1)
	if p.Quantity != 200 {
		t.Errorf("UpdateQuantity() quantity = %d, want 200", p.Quantity)
	}
}

func TestProductService_UpdateQuantity_NotFound(t *testing.T) {
	svc := NewProductService()

	err := svc.UpdateQuantity(999, 10)
	if err != ErrProductNotFound {
		t.Errorf("UpdateQuantity() error = %v, want %v", err, ErrProductNotFound)
	}
}

func TestProductService_UpdateQuantity_Invalid(t *testing.T) {
	svc := NewProductService()
	svc.AddProduct("Apple", 1.50, 100)

	err := svc.UpdateQuantity(1, -5)
	if err != ErrInvalidProduct {
		t.Errorf("UpdateQuantity() error = %v, want %v", err, ErrInvalidProduct)
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	svc := NewProductService()
	svc.AddProduct("Apple", 1.50, 100)

	err := svc.DeleteProduct(1)
	if err != nil {
		t.Fatalf("DeleteProduct() unexpected error: %v", err)
	}

	_, err = svc.GetProduct(1)
	if err != ErrProductNotFound {
		t.Errorf("GetProduct() after delete error = %v, want %v", err, ErrProductNotFound)
	}
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	svc := NewProductService()

	err := svc.DeleteProduct(999)
	if err != ErrProductNotFound {
		t.Errorf("DeleteProduct() error = %v, want %v", err, ErrProductNotFound)
	}
}

// TestProductService_ConcurrentAccess validates thread-safety for the race detector.
func TestProductService_ConcurrentAccess(t *testing.T) {
	svc := NewProductService()
	var wg sync.WaitGroup

	// Concurrently add products
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc.AddProduct("ConcurrentProduct", 9.99, 10)
		}()
	}

	// Concurrently read products
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc.ListProducts()
		}()
	}

	wg.Wait()

	products := svc.ListProducts()
	if len(products) != 50 {
		t.Errorf("ConcurrentAccess() product count = %d, want 50", len(products))
	}
}

// TestProductService_ConcurrentReadWrite tests concurrent reads and writes.
func TestProductService_ConcurrentReadWrite(t *testing.T) {
	svc := NewProductService()
	id, _ := svc.AddProduct("SharedProduct", 5.00, 100)

	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(qty int) {
			defer wg.Done()
			svc.UpdateQuantity(id, qty)
		}(i)
	}

	// Concurrent readers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc.GetProduct(id)
		}()
	}

	wg.Wait()
}
