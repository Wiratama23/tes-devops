package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/shopspring/decimal"

	"rwiratama.com/m/internal/handlers"
	"rwiratama.com/m/internal/models"
)

const (
	productInsertSQL = `
		INSERT INTO products (product_id, product_name, product_quantity, product_prices, product_type, created_by, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
	`
	productSelectByIDSQL = `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		WHERE product_id = $1
	`
	productUpdateSQL = `
		UPDATE products
		SET product_name = $1, product_quantity = $2, product_prices = $3, product_type = $4, image_path = $5
		WHERE product_id = $6
		RETURNING product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
	`
	productDeleteSQL = `DELETE FROM products WHERE product_id = $1`

	productPaginatedSQL = `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
)

func productRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"product_id", "product_name", "product_quantity", "product_prices",
		"product_type", "created_at", "created_by", "image_path",
	})
}

// CreateProduct -------------------------------------------------------------

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	mock := newMockPool(t)
	createdBy := uuid.New()
	now := time.Now().UTC()
	price := decimal.RequireFromString("29.99")

	mock.ExpectQuery(productInsertSQL).
		WithArgs("SKU1", "Widget", 100, pgxmock.AnyArg(), "10", createdBy, "assets/test.jpg").
		WillReturnRows(productRows().AddRow("SKU1", "Widget", 100, price, "10", now, createdBy, "assets/test.jpg"))

	h := handlers.NewProductHandler(mock)
	body := bytes.NewBufferString(fmt.Sprintf(
		`{"product_id":"SKU1","product_name":"Widget","product_quantity":100,"product_prices":"29.99","product_type":"10","created_by":"%s","image_path":"assets/test.jpg"}`,
		createdBy,
	))
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	h.CreateProduct(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.Product
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ProductID != "SKU1" || got.ProductName != "Widget" || got.ProductQuantity != 100 {
		t.Errorf("unexpected product: %+v", got)
	}
	if !got.ProductPrices.Equal(price) {
		t.Errorf("expected price %s, got %s", price, got.ProductPrices)
	}
	if got.ImagePath != "assets/test.jpg" {
		t.Errorf("expected image_path 'assets/test.jpg', got %q", got.ImagePath)
	}
}

func TestProductHandler_CreateProduct_DefaultImagePath(t *testing.T) {
	mock := newMockPool(t)
	createdBy := uuid.New()
	now := time.Now().UTC()
	price := decimal.RequireFromString("19.99")

	// When ImagePath is empty in the request, the handler should substitute the
	// default before passing to the repo, so we expect the default value here.
	mock.ExpectQuery(productInsertSQL).
		WithArgs("SKU2", "NoImage", 50, pgxmock.AnyArg(), "05", createdBy, "assets/default_image.jpg").
		WillReturnRows(productRows().AddRow("SKU2", "NoImage", 50, price, "05", now, createdBy, "assets/default_image.jpg"))

	h := handlers.NewProductHandler(mock)
	body := bytes.NewBufferString(fmt.Sprintf(
		`{"product_id":"SKU2","product_name":"NoImage","product_quantity":50,"product_prices":"19.99","product_type":"05","created_by":"%s"}`,
		createdBy,
	))
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	h.CreateProduct(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.Product
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ImagePath != "assets/default_image.jpg" {
		t.Errorf("expected default image_path, got %q", got.ImagePath)
	}
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewProductHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader("not-json"))
	w := httptest.NewRecorder()

	h.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestProductHandler_CreateProduct_RepoError(t *testing.T) {
	mock := newMockPool(t)
	createdBy := uuid.New()

	mock.ExpectQuery(productInsertSQL).
		WithArgs("SKU3", "Boom", 1, pgxmock.AnyArg(), "10", createdBy, "assets/x.jpg").
		WillReturnError(errors.New("db down"))

	h := handlers.NewProductHandler(mock)
	body := bytes.NewBufferString(fmt.Sprintf(
		`{"product_id":"SKU3","product_name":"Boom","product_quantity":1,"product_prices":"1","product_type":"10","created_by":"%s","image_path":"assets/x.jpg"}`,
		createdBy,
	))
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	h.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// GetProductByID -----------------------------------------------------------

func TestProductHandler_GetProductByID_Success(t *testing.T) {
	mock := newMockPool(t)
	createdBy := uuid.New()
	now := time.Now().UTC()
	price := decimal.RequireFromString("12.50")

	mock.ExpectQuery(productSelectByIDSQL).
		WithArgs("SKU1").
		WillReturnRows(productRows().AddRow("SKU1", "Widget", 7, price, "10", now, createdBy, "assets/x.jpg"))

	h := handlers.NewProductHandler(mock)
	req := chiRequest(http.MethodGet, "/products/SKU1", nil, map[string]string{"id": "SKU1"})
	w := httptest.NewRecorder()

	h.GetProductByID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.Product
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ProductID != "SKU1" {
		t.Errorf("expected SKU1, got %s", got.ProductID)
	}
	if !got.ProductPrices.Equal(price) {
		t.Errorf("expected price %s, got %s", price, got.ProductPrices)
	}
}

func TestProductHandler_GetProductByID_NotFound(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(productSelectByIDSQL).
		WithArgs("missing").
		WillReturnError(errors.New("no rows"))

	h := handlers.NewProductHandler(mock)
	req := chiRequest(http.MethodGet, "/products/missing", nil, map[string]string{"id": "missing"})
	w := httptest.NewRecorder()

	h.GetProductByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// GetAllProducts -----------------------------------------------------------

func TestProductHandler_GetAllProducts_PaginationShape(t *testing.T) {
	mock := newMockPool(t)
	createdBy := uuid.New()
	now := time.Now().UTC()
	price := decimal.RequireFromString("9.99")

	// page=3 limit=4 -> offset=8
	mock.ExpectQuery(productPaginatedSQL).
		WithArgs(4, 8).
		WillReturnRows(productRows().AddRow("SKU1", "A", 1, price, "10", now, createdBy, "assets/a.jpg"))

	h := handlers.NewProductHandler(mock)
	req := withPagination(httptest.NewRequest(http.MethodGet, "/products", nil), 3, 4)
	w := httptest.NewRecorder()

	h.GetAllProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got handlers.PaginatedProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Limit != 4 || got.Offset != 8 {
		t.Errorf("expected limit=4 offset=8, got limit=%d offset=%d", got.Limit, got.Offset)
	}
	if len(got.Data) != 1 {
		t.Errorf("expected 1 product, got %d", len(got.Data))
	}
}

func TestProductHandler_GetAllProducts_DefaultPagination(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(productPaginatedSQL).
		WithArgs(10, 0).
		WillReturnRows(productRows())

	h := handlers.NewProductHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	h.GetAllProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got handlers.PaginatedProductsResponse
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Limit != 10 || got.Offset != 0 {
		t.Errorf("expected default pagination, got limit=%d offset=%d", got.Limit, got.Offset)
	}
}

func TestProductHandler_GetAllProducts_RepoError(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(productPaginatedSQL).
		WithArgs(10, 0).
		WillReturnError(errors.New("db down"))

	h := handlers.NewProductHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	h.GetAllProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// UpdateProduct ------------------------------------------------------------

func TestProductHandler_UpdateProduct_Success(t *testing.T) {
	mock := newMockPool(t)
	createdBy := uuid.New()
	now := time.Now().UTC()
	price := decimal.RequireFromString("39.99")

	mock.ExpectQuery(productUpdateSQL).
		WithArgs("Widget2", 200, pgxmock.AnyArg(), "05", "assets/new.jpg", "SKU1").
		WillReturnRows(productRows().AddRow("SKU1", "Widget2", 200, price, "05", now, createdBy, "assets/new.jpg"))

	h := handlers.NewProductHandler(mock)
	body := bytes.NewBufferString(`{"product_name":"Widget2","product_quantity":200,"product_prices":"39.99","product_type":"05","image_path":"assets/new.jpg"}`)
	req := chiRequest(http.MethodPut, "/products/SKU1", body, map[string]string{"id": "SKU1"})
	w := httptest.NewRecorder()

	h.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.Product
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ProductName != "Widget2" || got.ProductQuantity != 200 {
		t.Errorf("unexpected product: %+v", got)
	}
	if !got.ProductPrices.Equal(price) {
		t.Errorf("expected price %s, got %s", price, got.ProductPrices)
	}
}

func TestProductHandler_UpdateProduct_InvalidJSON(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewProductHandler(mock)
	req := chiRequest(http.MethodPut, "/products/SKU1", strings.NewReader("not-json"), map[string]string{"id": "SKU1"})
	w := httptest.NewRecorder()

	h.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestProductHandler_UpdateProduct_RepoError(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(productUpdateSQL).
		WithArgs("Widget2", 200, pgxmock.AnyArg(), "05", "assets/new.jpg", "SKU1").
		WillReturnError(errors.New("db down"))

	h := handlers.NewProductHandler(mock)
	body := bytes.NewBufferString(`{"product_name":"Widget2","product_quantity":200,"product_prices":"39.99","product_type":"05","image_path":"assets/new.jpg"}`)
	req := chiRequest(http.MethodPut, "/products/SKU1", body, map[string]string{"id": "SKU1"})
	w := httptest.NewRecorder()

	h.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// DeleteProduct ------------------------------------------------------------

func TestProductHandler_DeleteProduct_Success(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectExec(productDeleteSQL).
		WithArgs("SKU1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := handlers.NewProductHandler(mock)
	req := chiRequest(http.MethodDelete, "/products/SKU1", nil, map[string]string{"id": "SKU1"})
	w := httptest.NewRecorder()

	h.DeleteProduct(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestProductHandler_DeleteProduct_NotFound(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectExec(productDeleteSQL).
		WithArgs("missing").
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	h := handlers.NewProductHandler(mock)
	req := chiRequest(http.MethodDelete, "/products/missing", nil, map[string]string{"id": "missing"})
	w := httptest.NewRecorder()

	h.DeleteProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
