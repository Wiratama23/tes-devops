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
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"

	"rwiratama.com/m/internal/handlers"
	"rwiratama.com/m/internal/models"
)

const (
	articleInsertSQL = `
		INSERT INTO articles (uid, title, article_text)
		VALUES ($1, $2, $3)
		RETURNING articles_id, uid, title, article_text, date_created, updated_at
	`
	articleSelectByIDSQL = `
		SELECT articles_id, uid, title, article_text, date_created, updated_at
		FROM articles
		WHERE articles_id = $1
	`
	articleSelectByUIDSQL = `
		SELECT articles_id, uid, title, article_text, date_created, updated_at
		FROM articles
		WHERE uid = $1
		ORDER BY date_created DESC
	`
	articleUpdateSQL = `
		UPDATE articles
		SET title = $2, article_text = $3, updated_at = CURRENT_TIMESTAMP
		WHERE articles_id = $1
		RETURNING articles_id, uid, title, article_text, date_created, updated_at
	`
	articleDeleteSQL = `DELETE FROM articles WHERE articles_id = $1`

	articlePaginatedSQL = `
		SELECT articles_id, uid, title, article_text, date_created, updated_at
		FROM articles
		ORDER BY date_created DESC
		LIMIT $1 OFFSET $2
	`
)

func articleRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{"articles_id", "uid", "title", "article_text", "date_created", "updated_at"})
}

// CreateArticle -------------------------------------------------------------

func TestArticleHandler_CreateArticle_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(articleInsertSQL).
		WithArgs(uid, "Hello", "World").
		WillReturnRows(articleRows().AddRow(42, uid, "Hello", "World", now, now))

	h := handlers.NewArticleHandler(mock)
	body := bytes.NewBufferString(fmt.Sprintf(`{"uid":"%s","title":"Hello","article_text":"World"}`, uid))
	req := httptest.NewRequest(http.MethodPost, "/articles", body)
	w := httptest.NewRecorder()

	h.CreateArticle(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.Article
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ArticlesID != 42 || got.Title != "Hello" || got.ArticleText != "World" {
		t.Errorf("unexpected article: %+v", got)
	}
}

func TestArticleHandler_CreateArticle_InvalidJSON(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewArticleHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/articles", strings.NewReader("not-json"))
	w := httptest.NewRecorder()

	h.CreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestArticleHandler_CreateArticle_RepoError(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()

	mock.ExpectQuery(articleInsertSQL).
		WithArgs(uid, "T", "B").
		WillReturnError(errors.New("boom"))

	h := handlers.NewArticleHandler(mock)
	body := bytes.NewBufferString(fmt.Sprintf(`{"uid":"%s","title":"T","article_text":"B"}`, uid))
	req := httptest.NewRequest(http.MethodPost, "/articles", body)
	w := httptest.NewRecorder()

	h.CreateArticle(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// GetArticle ----------------------------------------------------------------

func TestArticleHandler_GetArticle_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(articleSelectByIDSQL).
		WithArgs(7).
		WillReturnRows(articleRows().AddRow(7, uid, "T", "B", now, now))

	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodGet, "/articles/7", nil, map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	h.GetArticle(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.Article
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ArticlesID != 7 {
		t.Errorf("expected id 7, got %d", got.ArticlesID)
	}
}

func TestArticleHandler_GetArticle_InvalidID(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodGet, "/articles/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	h.GetArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestArticleHandler_GetArticle_NotFound(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(articleSelectByIDSQL).
		WithArgs(99).
		WillReturnError(pgx.ErrNoRows)

	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodGet, "/articles/99", nil, map[string]string{"id": "99"})
	w := httptest.NewRecorder()

	h.GetArticle(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// GetAllArticles ------------------------------------------------------------

func TestArticleHandler_GetAllArticles_PaginationShape(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	// WithArgs(limit, offset) should match the pagination logic in the handler, which
	// Default limit=10, page=2 -> offset=10*(2-1)=10
	mock.ExpectQuery(articlePaginatedSQL).
		WithArgs(10, 10).
		WillReturnRows(articleRows().AddRow(1, uid, "A", "B", now, now))

	h := handlers.NewArticleHandler(mock)
	req := withPagination(httptest.NewRequest(http.MethodGet, "/articles", nil), 2) // page=2
	w := httptest.NewRecorder()

	h.GetAllArticles(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got handlers.PaginatedArticlesResponse
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Offset != 10 {
		t.Errorf("expected offset=10, got offset=%d", got.Offset)
	}
	if len(got.Data) != 1 {
		t.Errorf("expected 1 article, got %d", len(got.Data))
	}
}

func TestArticleHandler_GetAllArticles_EmptyReturnsArray(t *testing.T) {
	mock := newMockPool(t)

	// No pagination context set on request -> handler falls back to page=1, limit=10
	mock.ExpectQuery(articlePaginatedSQL).
		WithArgs(10, 0).
		WillReturnRows(articleRows())

	h := handlers.NewArticleHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	h.GetAllArticles(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got handlers.PaginatedArticlesResponse
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Data == nil {
		t.Errorf("expected empty slice, got nil")
	}
	if len(got.Data) != 0 {
		t.Errorf("expected 0 articles, got %d", len(got.Data))
	}
	if got.Limit != 10 || got.Offset != 0 {
		t.Errorf("expected default pagination, got limit=%d offset=%d", got.Limit, got.Offset)
	}
}

func TestArticleHandler_GetAllArticles_RepoError(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(articlePaginatedSQL).
		WithArgs(10, 0).
		WillReturnError(errors.New("db down"))

	h := handlers.NewArticleHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	h.GetAllArticles(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// GetUserArticles -----------------------------------------------------------

func TestArticleHandler_GetUserArticles_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(articleSelectByUIDSQL).
		WithArgs(uid).
		WillReturnRows(articleRows().
			AddRow(1, uid, "A", "B", now, now).
			AddRow(2, uid, "C", "D", now, now))

	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodGet, "/users/"+uid.String()+"/articles", nil, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.GetUserArticles(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got []models.Article
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 articles, got %d", len(got))
	}
}

func TestArticleHandler_GetUserArticles_InvalidUUID(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewArticleHandler(mock)

	req := chiRequest(http.MethodGet, "/users/bad/articles", nil, map[string]string{"uid": "bad"})
	w := httptest.NewRecorder()

	h.GetUserArticles(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestArticleHandler_GetUserArticles_EmptyReturnsArray(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()

	mock.ExpectQuery(articleSelectByUIDSQL).
		WithArgs(uid).
		WillReturnRows(articleRows())

	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodGet, "/users/"+uid.String()+"/articles", nil, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.GetUserArticles(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	if body != "[]" {
		t.Errorf("expected empty JSON array, got %q", body)
	}
}

// UpdateArticle -------------------------------------------------------------

func TestArticleHandler_UpdateArticle_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(articleUpdateSQL).
		WithArgs(7, "NewT", "NewB").
		WillReturnRows(articleRows().AddRow(7, uid, "NewT", "NewB", now, now))

	h := handlers.NewArticleHandler(mock)
	body := bytes.NewBufferString(`{"title":"NewT","article_text":"NewB"}`)
	req := chiRequest(http.MethodPut, "/articles/7", body, map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	h.UpdateArticle(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.Article
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Title != "NewT" || got.ArticleText != "NewB" {
		t.Errorf("unexpected article: %+v", got)
	}
}

func TestArticleHandler_UpdateArticle_InvalidID(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewArticleHandler(mock)
	body := bytes.NewBufferString(`{"title":"x","article_text":"y"}`)
	req := chiRequest(http.MethodPut, "/articles/abc", body, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	h.UpdateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestArticleHandler_UpdateArticle_InvalidJSON(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodPut, "/articles/7", strings.NewReader("not-json"), map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	h.UpdateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestArticleHandler_UpdateArticle_NotFound(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(articleUpdateSQL).
		WithArgs(99, "x", "y").
		WillReturnError(pgx.ErrNoRows)

	h := handlers.NewArticleHandler(mock)
	body := bytes.NewBufferString(`{"title":"x","article_text":"y"}`)
	req := chiRequest(http.MethodPut, "/articles/99", body, map[string]string{"id": "99"})
	w := httptest.NewRecorder()

	h.UpdateArticle(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// DeleteArticle -------------------------------------------------------------

func TestArticleHandler_DeleteArticle_Success(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectExec(articleDeleteSQL).
		WithArgs(7).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodDelete, "/articles/7", nil, map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	h.DeleteArticle(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestArticleHandler_DeleteArticle_InvalidID(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodDelete, "/articles/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	h.DeleteArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestArticleHandler_DeleteArticle_NotFound(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectExec(articleDeleteSQL).
		WithArgs(99).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	h := handlers.NewArticleHandler(mock)
	req := chiRequest(http.MethodDelete, "/articles/99", nil, map[string]string{"id": "99"})
	w := httptest.NewRecorder()

	h.DeleteArticle(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
