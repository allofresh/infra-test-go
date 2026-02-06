package model

import (
	"testing"
)

func TestProduct_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected bool
	}{
		{
			name:     "available when quantity is positive",
			product:  Product{ID: 1, Name: "Apple", Price: 1.50, Quantity: 10},
			expected: true,
		},
		{
			name:     "unavailable when quantity is zero",
			product:  Product{ID: 2, Name: "Banana", Price: 2.00, Quantity: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.product.IsAvailable()
			if got != tt.expected {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_TotalValue(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected float64
	}{
		{
			name:     "calculates total value",
			product:  Product{Price: 10.00, Quantity: 5},
			expected: 50.00,
		},
		{
			name:     "zero quantity returns zero",
			product:  Product{Price: 10.00, Quantity: 0},
			expected: 0.00,
		},
		{
			name:     "zero price returns zero",
			product:  Product{Price: 0.00, Quantity: 100},
			expected: 0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.product.TotalValue()
			if got != tt.expected {
				t.Errorf("TotalValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_ApplyDiscount(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		percent  float64
		expected float64
	}{
		{
			name:     "10 percent discount",
			product:  Product{Price: 100.00},
			percent:  10,
			expected: 90.00,
		},
		{
			name:     "50 percent discount",
			product:  Product{Price: 200.00},
			percent:  50,
			expected: 100.00,
		},
		{
			name:     "0 percent discount",
			product:  Product{Price: 100.00},
			percent:  0,
			expected: 100.00,
		},
		{
			name:     "100 percent discount",
			product:  Product{Price: 100.00},
			percent:  100,
			expected: 0.00,
		},
		{
			name:     "negative percent returns original price",
			product:  Product{Price: 100.00},
			percent:  -10,
			expected: 100.00,
		},
		{
			name:     "over 100 percent returns original price",
			product:  Product{Price: 100.00},
			percent:  150,
			expected: 100.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.product.ApplyDiscount(tt.percent)
			if got != tt.expected {
				t.Errorf("ApplyDiscount(%v) = %v, want %v", tt.percent, got, tt.expected)
			}
		})
	}
}
