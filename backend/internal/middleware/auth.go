package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// AuthCookieName is the cookie that carries the JWT issued by the auth handler.
const AuthCookieName = "auth_token"

// TokenFromAuthCookie retrieves the JWT from the `auth_token` cookie set by the
// auth handler. It returns "" if the cookie is missing.
func TokenFromAuthCookie(r *http.Request) string {
	cookie, err := r.Cookie(AuthCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// JWTVerifier wraps jwtauth.Verify with our supported token sources:
//   - Authorization: Bearer header
//   - auth_token cookie (set by /api/auth/login)
//   - the default jwt cookie (kept for backwards compatibility)
//
// Use this together with jwtauth.Authenticator on protected routes.
func JWTVerifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return jwtauth.Verify(ja,
		jwtauth.TokenFromHeader,
		TokenFromAuthCookie,
		jwtauth.TokenFromCookie,
	)
}
