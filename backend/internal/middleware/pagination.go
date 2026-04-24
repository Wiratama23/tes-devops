package middleware

import (
	"context"
	"net/http"
	"strconv"
)

// Define a custom type for context keys
type contextKey string
const paginationKey contextKey = "pagination"

type PaginationData struct {
	Page  int
	Limit int
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get query parameters from the URL (e.g., /articles?page=2&limit=20)
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		// 2. Set default values
		page := 1
		limit := 10

		// 3. Convert strings to integers (ignoring errors for simplicity here)
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}

		// 4. Package the data
		data := PaginationData{Page: page, Limit: limit}

		// 5. Store it in the request context
		ctx := context.WithValue(r.Context(), paginationKey, data)

		// 6. Pass the request to the next handler (e.g., listArticles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPaginationData(ctx context.Context) PaginationData {
	if data, ok := ctx.Value(paginationKey).(PaginationData); ok {
		return data
	}
	// Fallback in case the middleware was forgotten on the route
	return PaginationData{Page: 1, Limit: 10} 
}