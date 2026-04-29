package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"rwiratama.com/m/internal/repository"
)

// AuthHandler issues and validates JWTs for the API. The cookie is set on the
// `/api` path so the browser sends it on every protected call. The same JWT
// can also be used as a `Authorization: Bearer <token>` header by SSR fetches.
type AuthHandler struct {
	repo      *repository.UserRepository
	tokenAuth *jwtauth.JWTAuth
	tokenTTL  time.Duration
	cookieName string
	secure     bool
}

type AuthHandlerConfig struct {
	TokenTTL   time.Duration
	CookieName string
	Secure     bool
}

func NewAuthHandler(pool repository.PgxPool, tokenAuth *jwtauth.JWTAuth, cfg AuthHandlerConfig) *AuthHandler {
	if cfg.TokenTTL == 0 {
		cfg.TokenTTL = 24 * time.Hour
	}
	if cfg.CookieName == "" {
		cfg.CookieName = "auth_token"
	}
	return &AuthHandler{
		repo:       repository.NewUserRepository(pool),
		tokenAuth:  tokenAuth,
		tokenTTL:   cfg.TokenTTL,
		cookieName: cfg.CookieName,
		secure:     cfg.Secure,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string  `json:"token"`
	User    UserDTO `json:"user"`
	Expires int64   `json:"expires"`
}

type UserDTO struct {
	UID      uuid.UUID `json:"uid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	IsAdmin  bool      `json:"is_admin"`
}

// Login verifies username/password against bcrypt and issues a JWT cookie.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	creds, err := h.repo.GetCredentialsByUsername(r.Context(), req.Username)
	if err != nil {
		// Avoid leaking which side failed.
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if creds.PasswordHash == "" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	expires := time.Now().Add(h.tokenTTL)
	claims := map[string]any{
		"uid":      creds.UID.String(),
		"username": creds.Username,
		"is_admin": creds.IsAdmin,
		"exp":      expires.Unix(),
	}

	_, tokenString, err := h.tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "failed to issue token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})

	resp := LoginResponse{
		Token: tokenString,
		User: UserDTO{
			UID:      creds.UID,
			Username: creds.Username,
			Email:    creds.Email,
			IsAdmin:  creds.IsAdmin,
		},
		Expires: expires.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Logout clears the auth cookie. It does not enforce that the caller is
// authenticated — clearing a cookie for an anonymous browser is a no-op.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}

// Me returns the current user resolved from the JWT in the request context.
// This route must be wrapped with jwtauth.Verifier + jwtauth.Authenticator.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uidRaw, _ := claims["uid"].(string)
	uid, err := uuid.Parse(uidRaw)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	user, err := h.repo.GetByID(r.Context(), uid)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	dto := UserDTO{
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
