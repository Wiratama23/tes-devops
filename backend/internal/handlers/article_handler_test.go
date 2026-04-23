package handlers

import (
	"bytes"
	"context"
	json "github.com/goccy/go-json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"rwiratama.com/m/internal/models"
	"rwiratama.com/m/internal/repository"
)

func TestArticleHandlerCreateArticle(t *testing.T) {
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

	handler := NewArticleHandler(pool)

	reqBody := CreateArticleRequest{
		UID:         user.UID,
		Title:       "Test Article",
		ArticleText: "This is a test article",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateArticle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var result models.Article
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Title != "Test Article" || result.ArticleText != "This is a test article" {
		t.Errorf("Article data mismatch: got %v", result)
	}
}

func TestArticleHandlerGetAllArticles(t *testing.T) {
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

	handler := NewArticleHandler(pool)

	req := httptest.NewRequest("GET", "/articles", nil)
	w := httptest.NewRecorder()

	handler.GetAllArticles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result []*models.Article
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result == nil {
		t.Error("Expected articles list, got nil")
	}
}

func TestArticleHandlerGetArticle(t *testing.T) {
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

	// Create a user and article first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	articleRepo := repository.NewArticleRepository(pool)
	article, _ := articleRepo.Create(ctx, user.UID, "Test Article", "This is a test article")

	// Get the article
	handler := NewArticleHandler(pool)
	req := httptest.NewRequest("GET", fmt.Sprintf("/articles/%d", article.ArticlesID), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", article.ArticlesID))
	w := httptest.NewRecorder()

	handler.GetArticle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.Article
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.ArticlesID != article.ArticlesID {
		t.Errorf("Expected article %d, got %d", article.ArticlesID, result.ArticlesID)
	}
}

func TestArticleHandlerGetUserArticles(t *testing.T) {
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

	// Create a user and articles
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	articleRepo := repository.NewArticleRepository(pool)
	if _, err := articleRepo.Create(ctx, user.UID, "Article 1", "Content 1"); err != nil {
		t.Fatalf("failed to create article 1: %v", err)
	}
	if _, err := articleRepo.Create(ctx, user.UID, "Article 2", "Content 2"); err != nil {
		t.Fatalf("failed to create article 2: %v", err)
	}

	// Get user articles
	handler := NewArticleHandler(pool)
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s/articles", user.UID.String()), nil)
	req.SetPathValue("uid", user.UID.String())
	w := httptest.NewRecorder()

	handler.GetUserArticles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result []*models.Article
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(result))
	}
}

func TestArticleHandlerUpdateArticle(t *testing.T) {
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

	// Create a user and article first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	articleRepo := repository.NewArticleRepository(pool)
	article, _ := articleRepo.Create(ctx, user.UID, "Original Title", "Original content")

	// Update the article
	handler := NewArticleHandler(pool)
	updateReq := UpdateArticleRequest{
		Title:       "Updated Title",
		ArticleText: "Updated content",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/articles/%d", article.ArticlesID), bytes.NewReader(body))
	req.SetPathValue("id", fmt.Sprintf("%d", article.ArticlesID))
	w := httptest.NewRecorder()

	handler.UpdateArticle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.Article
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Title != "Updated Title" || result.ArticleText != "Updated content" {
		t.Errorf("Article update failed: got %v", result)
	}
}

func TestArticleHandlerDeleteArticle(t *testing.T) {
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

	// Create a user and article first
	userRepo := repository.NewUserRepository(pool)
	user, _ := userRepo.Create(ctx, "testuser", "test@example.com")

	articleRepo := repository.NewArticleRepository(pool)
	article, _ := articleRepo.Create(ctx, user.UID, "Article to Delete", "Content to delete")

	// Delete the article
	handler := NewArticleHandler(pool)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/articles/%d", article.ArticlesID), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", article.ArticlesID))
	w := httptest.NewRecorder()

	handler.DeleteArticle(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify deletion
	_, err = articleRepo.GetByID(ctx, article.ArticlesID)
	if err == nil {
		t.Error("Expected article to be deleted")
	}
}

func TestArticleHandlerInvalidID(t *testing.T) {
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

	handler := NewArticleHandler(pool)
	req := httptest.NewRequest("GET", "/articles/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestArticleHandlerInvalidUID(t *testing.T) {
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

	handler := NewArticleHandler(pool)
	req := httptest.NewRequest("GET", "/users/invalid-uid/articles", nil)
	req.SetPathValue("uid", "invalid-uid")
	w := httptest.NewRecorder()

	handler.GetUserArticles(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestArticleHandlerInvalidMethod(t *testing.T) {
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

	handler := NewArticleHandler(pool)
	req := httptest.NewRequest("DELETE", "/articles", nil)
	w := httptest.NewRecorder()

	handler.CreateArticle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
