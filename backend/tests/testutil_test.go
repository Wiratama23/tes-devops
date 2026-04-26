package tests

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock/v5"

	czm "rwiratama.com/m/internal/middleware"
)

// chiRequest builds an *http.Request whose context carries a chi route-context
// populated with the given URL params, so handlers that call
// chi.URLParam(r, name) resolve correctly without spinning up a full router.
func chiRequest(method, target string, body io.Reader, params map[string]string) *http.Request {
	req := httptest.NewRequest(method, target, body)
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

// withPagination wraps the given request through the real Paginate middleware
// so that middleware.GetPaginationData(r.Context()) returns the requested
// page/limit pair. The pagination context key is unexported, so reusing the
// middleware itself is the only way to populate it correctly.
func withPagination(r *http.Request, page, limit int) *http.Request {
	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	r.URL = u

	var captured *http.Request
	czm.Paginate(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		captured = req
	})).ServeHTTP(httptest.NewRecorder(), r)
	return captured
}

// newMockPool returns a pgxmock pool wired with QueryMatcherEqual so SQL
// strings must match exactly, plus a t.Cleanup hook that asserts every queued
// expectation was met and closes the pool.
// Use QueryMatcherRegexp to match SQL strings exactly, ignoring parameter order.
func newMockPool(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet pgxmock expectations: %v", err)
		}
		mock.Close()
	})
	return mock
}
