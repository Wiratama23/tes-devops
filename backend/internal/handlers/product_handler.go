package handlers

import (
	json "github.com/goccy/go-json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"rwiratama.com/m/internal/repository"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(pool *pgxpool.Pool) *ProductHandler {
	return &ProductHandler{
		repo: repository.NewProductRepository(pool),
	}
}

type CreateProductRequest struct {
	ProductID       string          `json:"product_id"`
	ProductName     string          `json:"product_name"`
	ProductQuantity int             `json:"product_quantity"`
	ProductPrices   decimal.Decimal `json:"product_prices"`
	ProductType     string          `json:"product_type"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	ImagePath       string          `json:"image_path"`
}

type UpdateProductRequest struct {
	ProductName     string          `json:"product_name"`
	ProductQuantity int             `json:"product_quantity"`
	ProductPrices   decimal.Decimal `json:"product_prices"`
	ProductType     string          `json:"product_type"`
	ImagePath       string          `json:"image_path"`
}

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set default image path if not provided
	if req.ImagePath == "" {
		req.ImagePath = "assets/default_image.jpg"
	}

	product, err := h.repo.Create(r.Context(), req.ProductID, req.ProductName, req.ProductQuantity, req.ProductPrices.String(), req.ProductType, req.CreatedBy, req.ImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetProductByID handles GET /products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	product, err := h.repo.GetByID(r.Context(), productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAllProducts handles GET /products
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateProduct handles PUT /products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.repo.Update(r.Context(), productID, req.ProductName, req.ProductQuantity, req.ProductPrices.String(), req.ProductType, req.ImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteProduct handles DELETE /products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	err := h.repo.Delete(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
