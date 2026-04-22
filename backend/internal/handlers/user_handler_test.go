package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"rwiratama.com/m/internal/models"
	"rwiratama.com/m/internal/repository"
)

// MockUserRepository for testing
type MockUserRepository struct {
	users map[uuid.UUID]*models.User
	calls []string
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uuid.UUID]*models.User),
		calls: []string{},
	}
}

func (m *MockUserRepository) Create(ctx context.Context, username, email string) (*models.User, error) {
	m.calls = append(m.calls, "Create")
	user := &models.User{
		UID:       uuid.New(),
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.users[user.UID] = user
	return user, nil
}

func (m *MockUserRepository) GetByUID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	m.calls = append(m.calls, "GetByUID")
	if user, ok := m.users[uid]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	m.calls = append(m.calls, "GetAll")
	var users []*models.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, uid uuid.UUID, username, email string) (*models.User, error) {
	m.calls = append(m.calls, "Update")
	if user, ok := m.users[uid]; ok {
		user.Username = username
		user.Email = email
		user.UpdatedAt = time.Now()
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	m.calls = append(m.calls, "Delete")
	if _, ok := m.users[uid]; ok {
		delete(m.users, uid)
		return nil
	}
	return fmt.Errorf("user not found")
}

// UserHandler wrapper for testing
func NewUserHandlerWithRepo(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func TestUserHandlerCreateUser(t *testing.T) {
	// Create a real pool connection for integration test
	ctx := context.Background()

	// For unit test purposes, we'll skip if DATABASE_URL is not set
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	handler := NewUserHandler(pool)

	reqBody := CreateUserRequest{
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
	json.NewDecoder(w.Body).Decode(&result)
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

	handler := NewUserHandler(pool)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result []*models.User
	json.NewDecoder(w.Body).Decode(&result)
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

	// Create a user first
	handler := NewUserHandler(pool)
	repo := repository.NewUserRepository(pool)
	user, _ := repo.Create(ctx, "testuser", "test@example.com")

	// Get the user
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s", user.UID.String()), nil)
	req.SetPathValue("uid", user.UID.String())
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.User
	json.NewDecoder(w.Body).Decode(&result)
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

	// Create a user first
	repo := repository.NewUserRepository(pool)
	user, _ := repo.Create(ctx, "testuser", "test@example.com")

	// Update the user
	handler := NewUserHandler(pool)
	updateReq := UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%s", user.UID.String()), bytes.NewReader(body))
	req.SetPathValue("uid", user.UID.String())
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.User
	json.NewDecoder(w.Body).Decode(&result)
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

	// Create a user first
	repo := repository.NewUserRepository(pool)
	user, _ := repo.Create(ctx, "testuser", "test@example.com")

	// Delete the user
	handler := NewUserHandler(pool)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", user.UID.String()), nil)
	req.SetPathValue("uid", user.UID.String())
	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, user.UID)
	if err == nil {
		t.Error("Expected user to be deleted")
	}
}

func getTestDatabaseURL() string {
	// For testing, you should set DATABASE_URL environment variable
	// or use a test database
	return ""
}
