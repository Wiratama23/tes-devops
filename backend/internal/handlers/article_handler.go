package handlers

import (
	json "github.com/goccy/go-json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"rwiratama.com/m/internal/models"
	"rwiratama.com/m/internal/repository"
)

type ArticleHandler struct {
	repo *repository.ArticleRepository
}

type PaginatedArticlesResponse struct {
	Data       []models.Article `json:"data"`
	TotalCount int              `json:"total_count"`
	Limit      int              `json:"limit"`
	Offset     int              `json:"offset"`
}

func NewArticleHandler(pool *pgxpool.Pool) *ArticleHandler {
	return &ArticleHandler{
		repo: repository.NewArticleRepository(pool),
	}
}

type CreateArticleRequest struct {
	UID         uuid.UUID `json:"uid"`
	Title       string    `json:"title"`
	ArticleText string    `json:"article_text"`
}

type UpdateArticleRequest struct {
	Title       string `json:"title"`
	ArticleText string `json:"article_text"`
}

// CreateArticle handles POST /articles
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	article, err := h.repo.Create(r.Context(), req.UID, req.Title, req.ArticleText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetArticle handles GET /articles/{id}
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

// GetAllArticles handles GET /articles with optional ?limit and ?offset query parameters
func (h *ArticleHandler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	// Parse limit and offset from query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	articles, totalCount, err := h.repo.GetAllWithPagination(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if articles == nil {
		articles = []models.Article{}
	}

	resp := PaginatedArticlesResponse{
		Data:       articles,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetUserArticles handles GET /users/{uid}/articles
func (h *ArticleHandler) GetUserArticles(w http.ResponseWriter, r *http.Request) {
	uidStr := chi.URLParam(r, "uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	articles, err := h.repo.GetByUserID(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if articles == nil {
		articles = []models.Article{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

// UpdateArticle handles PUT /articles/{id}
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var req UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	article, err := h.repo.Update(r.Context(), id, req.Title, req.ArticleText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}	
}

// DeleteArticle handles DELETE /articles/{id}
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
