package tests

import (
	"bytes"
	"encoding/json"
	"errors"
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
	userInsertSQL = `
		INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING uid, username, email, is_admin, created_at, updated_at
	`
	userSelectByIDSQL = `
		SELECT uid, username, email, is_admin, created_at, updated_at
		FROM users
		WHERE uid = $1
	`
	userSelectAllSQL = `
		SELECT uid, username, email, is_admin, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`
	userUpdateSQL = `
		UPDATE users
		SET username = $2, email = $3, updated_at = CURRENT_TIMESTAMP
		WHERE uid = $1
		RETURNING uid, username, email, is_admin, created_at, updated_at
	`
	userDeleteSQL = `DELETE FROM users WHERE uid = $1`
)

var userColumns = []string{"uid", "username", "email", "is_admin", "created_at", "updated_at"}

func userRows(uid uuid.UUID, username, email string, ts time.Time) *pgxmock.Rows {
	return pgxmock.NewRows(userColumns).
		AddRow(uid, username, email, false, ts, ts)
}

// CreateUser ----------------------------------------------------------------

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(userInsertSQL).
		WithArgs("alice", "alice@example.com").
		WillReturnRows(userRows(uid, "alice", "alice@example.com", now))

	h := handlers.NewUserHandler(mock)

	body := bytes.NewBufferString(`{"username":"alice","email":"alice@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.UID != uid || got.Username != "alice" || got.Email != "alice@example.com" {
		t.Errorf("unexpected user: %+v", got)
	}
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("not-json"))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_CreateUser_RepoError(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(userInsertSQL).
		WithArgs("bob", "bob@example.com").
		WillReturnError(errors.New("boom"))

	h := handlers.NewUserHandler(mock)
	body := bytes.NewBufferString(`{"username":"bob","email":"bob@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// GetUser -------------------------------------------------------------------

func TestUserHandler_GetUser_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(userSelectByIDSQL).
		WithArgs(uid).
		WillReturnRows(userRows(uid, "alice", "alice@example.com", now))

	h := handlers.NewUserHandler(mock)
	req := chiRequest(http.MethodGet, "/users/"+uid.String(), nil, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.UID != uid {
		t.Errorf("expected uid %s, got %s", uid, got.UID)
	}
}

func TestUserHandler_GetUser_InvalidUUID(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewUserHandler(mock)
	req := chiRequest(http.MethodGet, "/users/not-a-uuid", nil, map[string]string{"uid": "not-a-uuid"})
	w := httptest.NewRecorder()

	h.GetUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()

	mock.ExpectQuery(userSelectByIDSQL).
		WithArgs(uid).
		WillReturnError(pgx.ErrNoRows)

	h := handlers.NewUserHandler(mock)
	req := chiRequest(http.MethodGet, "/users/"+uid.String(), nil, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.GetUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// GetAllUsers ---------------------------------------------------------------

func TestUserHandler_GetAllUsers_WithRows(t *testing.T) {
	mock := newMockPool(t)
	uid1, uid2 := uuid.New(), uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(userSelectAllSQL).
		WillReturnRows(pgxmock.NewRows(userColumns).
			AddRow(uid1, "alice", "a@example.com", false, now, now).
			AddRow(uid2, "bob", "b@example.com", true, now, now))

	h := handlers.NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got []models.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 users, got %d", len(got))
	}
}

func TestUserHandler_GetAllUsers_EmptyReturnsArray(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(userSelectAllSQL).
		WillReturnRows(pgxmock.NewRows(userColumns))

	h := handlers.NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	if body != "[]" {
		t.Errorf("expected empty JSON array, got %q", body)
	}
}

func TestUserHandler_GetAllUsers_RepoError(t *testing.T) {
	mock := newMockPool(t)

	mock.ExpectQuery(userSelectAllSQL).
		WillReturnError(errors.New("db down"))

	h := handlers.NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// UpdateUser ----------------------------------------------------------------

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	now := time.Now().UTC()

	mock.ExpectQuery(userUpdateSQL).
		WithArgs(uid, "alice2", "alice2@example.com").
		WillReturnRows(userRows(uid, "alice2", "alice2@example.com", now))

	h := handlers.NewUserHandler(mock)
	body := bytes.NewBufferString(`{"username":"alice2","email":"alice2@example.com"}`)
	req := chiRequest(http.MethodPut, "/users/"+uid.String(), body, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var got models.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Username != "alice2" || got.Email != "alice2@example.com" {
		t.Errorf("unexpected user: %+v", got)
	}
}

func TestUserHandler_UpdateUser_InvalidUUID(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewUserHandler(mock)
	body := bytes.NewBufferString(`{"username":"x","email":"x@y.io"}`)
	req := chiRequest(http.MethodPut, "/users/bad", body, map[string]string{"uid": "bad"})
	w := httptest.NewRecorder()

	h.UpdateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_UpdateUser_InvalidJSON(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	h := handlers.NewUserHandler(mock)

	req := chiRequest(http.MethodPut, "/users/"+uid.String(), strings.NewReader("not-json"), map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.UpdateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_UpdateUser_NotFound(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()

	mock.ExpectQuery(userUpdateSQL).
		WithArgs(uid, "x", "x@y.io").
		WillReturnError(pgx.ErrNoRows)

	h := handlers.NewUserHandler(mock)
	body := bytes.NewBufferString(`{"username":"x","email":"x@y.io"}`)
	req := chiRequest(http.MethodPut, "/users/"+uid.String(), body, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.UpdateUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// DeleteUser ----------------------------------------------------------------

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()

	mock.ExpectExec(userDeleteSQL).
		WithArgs(uid).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	h := handlers.NewUserHandler(mock)
	req := chiRequest(http.MethodDelete, "/users/"+uid.String(), nil, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestUserHandler_DeleteUser_InvalidUUID(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewUserHandler(mock)
	req := chiRequest(http.MethodDelete, "/users/bad", nil, map[string]string{"uid": "bad"})
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()

	mock.ExpectExec(userDeleteSQL).
		WithArgs(uid).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	h := handlers.NewUserHandler(mock)
	req := chiRequest(http.MethodDelete, "/users/"+uid.String(), nil, map[string]string{"uid": uid.String()})
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
