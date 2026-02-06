package model

// Product represents a product in the AlloFresh catalog.
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// IsAvailable checks whether the product is in stock.
func (p *Product) IsAvailable() bool {
	return p.Quantity > 0
}

// TotalValue returns the total inventory value for this product.
func (p *Product) TotalValue() float64 {
	return p.Price * float64(p.Quantity)
}

// ApplyDiscount applies a percentage discount (0-100) and returns the discounted price.
func (p *Product) ApplyDiscount(percent float64) float64 {
	if percent < 0 || percent > 100 {
		return p.Price
	}
	return p.Price * (1 - percent/100)
}
