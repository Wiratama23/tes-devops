package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
	"golang.org/x/crypto/bcrypt"

	"rwiratama.com/m/internal/handlers"
	czm "rwiratama.com/m/internal/middleware"
)

const credsSelectByUsernameSQL = `
		SELECT uid, username, email, password_hash, is_admin
		FROM users
		WHERE username = $1
	`

func newTokenAuth() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte("test-secret"), nil)
}

func bcryptHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return string(hash)
}

func TestAuthHandler_Login_Success_SetsCookie(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	hash := bcryptHash(t, "secret123")

	mock.ExpectQuery(credsSelectByUsernameSQL).
		WithArgs("alice").
		WillReturnRows(pgxmock.NewRows([]string{"uid", "username", "email", "password_hash", "is_admin"}).
			AddRow(uid, "alice", "alice@example.com", hash, true))

	tokenAuth := newTokenAuth()
	h := handlers.NewAuthHandler(mock, tokenAuth, handlers.AuthHandlerConfig{
		TokenTTL:   time.Hour,
		CookieName: czm.AuthCookieName,
	})

	body := bytes.NewBufferString(`{"username":"alice","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}

	var resp handlers.LoginResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Token == "" {
		t.Errorf("expected non-empty token")
	}
	if !resp.User.IsAdmin {
		t.Errorf("expected is_admin=true on response")
	}

	cookies := w.Result().Cookies()
	var found *http.Cookie
	for _, c := range cookies {
		if c.Name == czm.AuthCookieName {
			found = c
			break
		}
	}
	if found == nil {
		t.Fatalf("expected %s cookie", czm.AuthCookieName)
	}
	if found != nil && (!found.HttpOnly || found.SameSite != http.SameSiteLaxMode) {
		t.Errorf("expected HttpOnly + SameSite=Lax cookie, got %+v", found)
	}
}

func TestAuthHandler_Login_BadJSON(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewAuthHandler(mock, newTokenAuth(), handlers.AuthHandlerConfig{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader("{"))
	w := httptest.NewRecorder()
	h.Login(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewAuthHandler(mock, newTokenAuth(), handlers.AuthHandlerConfig{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	h.Login(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	mock := newMockPool(t)
	mock.ExpectQuery(credsSelectByUsernameSQL).
		WithArgs("ghost").
		WillReturnError(pgx.ErrNoRows)

	h := handlers.NewAuthHandler(mock, newTokenAuth(), handlers.AuthHandlerConfig{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"ghost","password":"x"}`))
	w := httptest.NewRecorder()
	h.Login(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()
	hash := bcryptHash(t, "correct-password")

	mock.ExpectQuery(credsSelectByUsernameSQL).
		WithArgs("alice").
		WillReturnRows(pgxmock.NewRows([]string{"uid", "username", "email", "password_hash", "is_admin"}).
			AddRow(uid, "alice", "alice@example.com", hash, false))

	h := handlers.NewAuthHandler(mock, newTokenAuth(), handlers.AuthHandlerConfig{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"alice","password":"wrong"}`))
	w := httptest.NewRecorder()
	h.Login(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_EmptyHashRejected(t *testing.T) {
	mock := newMockPool(t)
	uid := uuid.New()

	mock.ExpectQuery(credsSelectByUsernameSQL).
		WithArgs("alice").
		WillReturnRows(pgxmock.NewRows([]string{"uid", "username", "email", "password_hash", "is_admin"}).
			AddRow(uid, "alice", "alice@example.com", "", false))

	h := handlers.NewAuthHandler(mock, newTokenAuth(), handlers.AuthHandlerConfig{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"alice","password":"anything"}`))
	w := httptest.NewRecorder()
	h.Login(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for empty hash, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_ClearsCookie(t *testing.T) {
	mock := newMockPool(t)
	h := handlers.NewAuthHandler(mock, newTokenAuth(), handlers.AuthHandlerConfig{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()
	h.Logout(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == czm.AuthCookieName {
			if c.MaxAge >= 0 {
				t.Errorf("expected cookie cleared (MaxAge < 0), got %d", c.MaxAge)
			}
			return
		}
	}
	t.Errorf("expected logout to set %s cookie", czm.AuthCookieName)
}
