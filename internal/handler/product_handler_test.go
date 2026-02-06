package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allofresh/infra-test-go/internal/service"
)

func setupHandler() (*ProductHandler, *http.ServeMux) {
	svc := service.NewProductService()
	h := NewProductHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, mux
}

func TestHealthCheck(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HealthCheck status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Errorf("HealthCheck status = %s, want ok", resp["status"])
	}
}

func TestAddProduct(t *testing.T) {
	_, mux := setupHandler()

	body := `{"name":"Apple","price":1.50,"quantity":100}`
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("AddProduct status = %d, want %d", w.Code, http.StatusCreated)
	}

	var resp map[string]int
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] != 1 {
		t.Errorf("AddProduct id = %d, want 1", resp["id"])
	}
}

func TestAddProduct_InvalidBody(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("AddProduct invalid body status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAddProduct_EmptyName(t *testing.T) {
	_, mux := setupHandler()

	body := `{"name":"","price":1.50,"quantity":100}`
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("AddProduct empty name status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestListProducts(t *testing.T) {
	_, mux := setupHandler()

	// Add a product first
	body := `{"name":"Apple","price":1.50,"quantity":100}`
	addReq := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	addReq.Header.Set("Content-Type", "application/json")
	addW := httptest.NewRecorder()
	mux.ServeHTTP(addW, addReq)

	// List products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ListProducts status = %d, want %d", w.Code, http.StatusOK)
	}

	var products []map[string]interface{}
	json.NewDecoder(w.Body).Decode(&products)
	if len(products) != 1 {
		t.Errorf("ListProducts count = %d, want 1", len(products))
	}
}

func TestGetProduct(t *testing.T) {
	_, mux := setupHandler()

	// Add a product first
	body := `{"name":"Apple","price":1.50,"quantity":100}`
	addReq := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	addReq.Header.Set("Content-Type", "application/json")
	addW := httptest.NewRecorder()
	mux.ServeHTTP(addW, addReq)

	// Get the product
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetProduct status = %d, want %d", w.Code, http.StatusOK)
	}

	var product map[string]interface{}
	json.NewDecoder(w.Body).Decode(&product)
	if product["name"] != "Apple" {
		t.Errorf("GetProduct name = %v, want Apple", product["name"])
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GetProduct not found status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetProduct_InvalidID(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("GetProduct invalid id status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteProduct(t *testing.T) {
	_, mux := setupHandler()

	// Add a product first
	body := `{"name":"Apple","price":1.50,"quantity":100}`
	addReq := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	addReq.Header.Set("Content-Type", "application/json")
	addW := httptest.NewRecorder()
	mux.ServeHTTP(addW, addReq)

	// Delete it
	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("DeleteProduct status = %d, want %d", w.Code, http.StatusNoContent)
	}

	// Verify it's gone
	getReq := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	getW := httptest.NewRecorder()
	mux.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("GetProduct after delete status = %d, want %d", getW.Code, http.StatusNotFound)
	}
}

func TestDeleteProduct_NotFound(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodDelete, "/products/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("DeleteProduct not found status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
