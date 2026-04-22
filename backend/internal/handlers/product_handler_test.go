package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"rwiratama.com/m/internal/models"
	"rwiratama.com/m/internal/repository"
)

func TestProductHandlerCreateProduct(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Create a user first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	handler := NewProductHandler(pool)

	reqBody := CreateProductRequest{
		ProductID:       "SKU10001",
		ProductName:     "Test Product",
		ProductQuantity: 100,
		ProductPrices:   decimal.NewFromFloat(29.99),
		ProductType:     "10",
		CreatedBy:       user.UID,
		ImagePath:       "assets/test.jpg",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var result models.Product
	json.NewDecoder(w.Body).Decode(&result)
	if result.ProductID != "SKU10001" || result.ProductName != "Test Product" {
		t.Errorf("Product data mismatch: got %v", result)
	}
}

func TestProductHandlerCreateProductWithDefaultImage(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Create a user first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	handler := NewProductHandler(pool)

	reqBody := CreateProductRequest{
		ProductID:       "SKU10002",
		ProductName:     "Product Without Image",
		ProductQuantity: 50,
		ProductPrices:   decimal.NewFromFloat(19.99),
		ProductType:     "05",
		CreatedBy:       user.UID,
		// ImagePath is empty, should default to assets/default_image.jpg
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var result models.Product
	json.NewDecoder(w.Body).Decode(&result)
	if result.ImagePath != "assets/default_image.jpg" {
		t.Errorf("Expected default image path, got %s", result.ImagePath)
	}
}

func TestProductHandlerGetAllProducts(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	handler := NewProductHandler(pool)

	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	handler.GetAllProducts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result []*models.Product
	json.NewDecoder(w.Body).Decode(&result)
	if result == nil {
		t.Error("Expected products list, got nil")
	}
}

func TestProductHandlerGetProductByID(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Create a user and product first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	productRepo := repository.NewProductRepository(pool)
	product, _ := productRepo.Create(ctx, "SKU10003", "Test Product", 100, "29.99", "10", user.UID, "assets/test.jpg")

	// Get the product
	handler := NewProductHandler(pool)
	req := httptest.NewRequest("GET", fmt.Sprintf("/products/%s", product.ProductID), nil)
	req.SetPathValue("id", product.ProductID)
	w := httptest.NewRecorder()

	handler.GetProductByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.Product
	json.NewDecoder(w.Body).Decode(&result)
	if result.ProductID != product.ProductID {
		t.Errorf("Expected product %s, got %s", product.ProductID, result.ProductID)
	}
}

func TestProductHandlerUpdateProduct(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Create a user and product first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	productRepo := repository.NewProductRepository(pool)
	product, _ := productRepo.Create(ctx, "SKU10004", "Original Product", 100, "29.99", "10", user.UID, "assets/test.jpg")

	// Update the product
	handler := NewProductHandler(pool)
	updateReq := UpdateProductRequest{
		ProductName:     "Updated Product",
		ProductQuantity: 200,
		ProductPrices:   decimal.NewFromFloat(39.99),
		ProductType:     "05",
		ImagePath:       "assets/updated.jpg",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/products/%s", product.ProductID), bytes.NewReader(body))
	req.SetPathValue("id", product.ProductID)
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.Product
	json.NewDecoder(w.Body).Decode(&result)
	if result.ProductName != "Updated Product" || result.ProductQuantity != 200 {
		t.Errorf("Product update failed: got %v", result)
	}
}

func TestProductHandlerDeleteProduct(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Create a user and product first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	productRepo := repository.NewProductRepository(pool)
	product, _ := productRepo.Create(ctx, "SKU10005", "Product to Delete", 100, "29.99", "10", user.UID, "assets/test.jpg")

	// Delete the product
	handler := NewProductHandler(pool)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/products/%s", product.ProductID), nil)
	req.SetPathValue("id", product.ProductID)
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify deletion
	_, err = productRepo.GetByID(ctx, product.ProductID)
	if err == nil {
		t.Error("Expected product to be deleted")
	}
}

func TestProductHandlerInvalidMethod(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	handler := NewProductHandler(pool)

	// Try DELETE on POST endpoint
	req := httptest.NewRequest("DELETE", "/products", nil)
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestProductHandlerInvalidPrice(t *testing.T) {
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Create a user first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	handler := NewProductHandler(pool)

	reqBody := CreateProductRequest{
		ProductID:       "SKU10006",
		ProductName:     "Invalid Price Product",
		ProductQuantity: 100,
		ProductPrices:   decimal.NewFromFloat(29.99),
		ProductType:     "10",
		CreatedBy:       user.UID,
		ImagePath:       "assets/test.jpg",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	// Should succeed with valid decimal
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}
