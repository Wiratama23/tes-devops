package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rwiratama.com/m/internal/handlers"
)

func TestLogHandler_Receive_Success(t *testing.T) {
	h := handlers.NewLogHandler()

	body := bytes.NewBufferString(`{"level":"error","message":"client crash","stack":"at foo","url":"/admin","meta":{"a":1}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/logs", body)
	req.Header.Set("User-Agent", "test-agent/1.0")
	w := httptest.NewRecorder()

	h.Receive(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestLogHandler_Receive_BadJSON(t *testing.T) {
	h := handlers.NewLogHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader("not-json"))
	w := httptest.NewRecorder()

	h.Receive(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLogHandler_Receive_DefaultsLevel(t *testing.T) {
	h := handlers.NewLogHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader(`{"message":"hi"}`))
	w := httptest.NewRecorder()
	h.Receive(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d (body=%s)", w.Code, w.Body.String())
	}
}

func TestLogHandler_Receive_EnforcesPayloadCap(t *testing.T) {
	h := handlers.NewLogHandler()
	huge := strings.Repeat("x", 2<<20) // 2 MiB > 1 MiB cap
	body := `{"level":"error","message":"` + huge + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.Receive(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 due to body cap, got %d", w.Code)
	}
}
