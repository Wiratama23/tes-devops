package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	json "github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"rwiratama.com/m/internal/handlers"
	"rwiratama.com/m/internal/models"
	"rwiratama.com/m/internal/repository"
)

// withChiParam attaches a chi route-context with named URL params so handlers
// that call chi.URLParam(r, name) can resolve them in tests.
func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestUserHandlerCreateUser(t *testing.T) {
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

	handler := handlers.NewUserHandler(pool)

	reqBody := handlers.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var result models.User
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Username != "testuser" || result.Email != "test@example.com" {
		t.Errorf("User data mismatch: got %v", result)
	}
}

func TestUserHandlerGetAllUsers(t *testing.T) {
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

	handler := handlers.NewUserHandler(pool)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result []models.User
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result == nil {
		t.Error("Expected users list, got nil")
	}
}

func TestUserHandlerGetUser(t *testing.T) {
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

	handler := handlers.NewUserHandler(pool)
	repo := repository.NewUserRepository(pool)
	user, _ := repo.Create(ctx, "testuser", "test@example.com")

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s", user.UID.String()), nil)
	req = withChiParam(req, "uid", user.UID.String())
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.User
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.UID != user.UID {
		t.Errorf("Expected user %s, got %s", user.UID, result.UID)
	}
}

func TestUserHandlerUpdateUser(t *testing.T) {
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

	repo := repository.NewUserRepository(pool)
	user, _ := repo.Create(ctx, "testuser", "test@example.com")

	handler := handlers.NewUserHandler(pool)
	updateReq := handlers.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%s", user.UID.String()), bytes.NewReader(body))
	req = withChiParam(req, "uid", user.UID.String())
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.User
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Username != "updateduser" || result.Email != "updated@example.com" {
		t.Errorf("User update failed: got %v", result)
	}
}

func TestUserHandlerDeleteUser(t *testing.T) {
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

	repo := repository.NewUserRepository(pool)
	user, _ := repo.Create(ctx, "testuser", "test@example.com")

	handler := handlers.NewUserHandler(pool)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", user.UID.String()), nil)
	req = withChiParam(req, "uid", user.UID.String())
	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	if _, err := repo.GetByID(ctx, user.UID); err == nil {
		t.Error("Expected user to be deleted")
	}
}
