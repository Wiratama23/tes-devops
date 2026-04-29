package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// RequireAdmin is a chi middleware that allows the request through only when
// the JWT carries `is_admin: true`. It must be installed *after* the standard
// jwtauth.Verifier + jwtauth.Authenticator pair so the claims are already
// populated on the context.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		isAdmin, _ := claims["is_admin"].(bool)
		if !isAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
