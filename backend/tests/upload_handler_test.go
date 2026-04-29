package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"rwiratama.com/m/internal/handlers"
)

func multipartImage(t *testing.T, fieldName, filename, contentType string, data []byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+filename+`"`)
	if contentType != "" {
		header.Set("Content-Type", contentType)
	}
	part, err := w.CreatePart(header)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	return &buf, w.FormDataContentType()
}

func TestUploadHandler_UploadImage_Success(t *testing.T) {
	dir := t.TempDir()
	h := handlers.NewUploadHandler(dir)

	body, ct := multipartImage(t, "file", "cat.png", "image/png", []byte("PNGDATA"))
	req := httptest.NewRequest(http.MethodPost, "/api/uploads/images", body)
	req.Header.Set("Content-Type", ct)
	w := httptest.NewRecorder()

	h.UploadImage(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%s)", w.Code, w.Body.String())
	}

	var resp handlers.UploadResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Filename == "" || resp.URL == "" || resp.Path == "" {
		t.Errorf("missing fields: %+v", resp)
	}
	if filepath.Ext(resp.Filename) != ".png" {
		t.Errorf("expected .png filename, got %s", resp.Filename)
	}
	if _, err := os.Stat(filepath.Join(dir, resp.Filename)); err != nil {
		t.Errorf("expected file on disk: %v", err)
	}
}

func TestUploadHandler_UploadImage_RejectsBadExtension(t *testing.T) {
	dir := t.TempDir()
	h := handlers.NewUploadHandler(dir)

	body, ct := multipartImage(t, "file", "evil.exe", "application/octet-stream", []byte("MZ"))
	req := httptest.NewRequest(http.MethodPost, "/api/uploads/images", body)
	req.Header.Set("Content-Type", ct)
	w := httptest.NewRecorder()

	h.UploadImage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUploadHandler_UploadImage_MissingField(t *testing.T) {
	dir := t.TempDir()
	h := handlers.NewUploadHandler(dir)

	var buf bytes.Buffer
	w2 := multipart.NewWriter(&buf)
	if err := w2.WriteField("not_file", "x"); err != nil {
		t.Fatalf("write field: %v", err)
	}
	w2.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/uploads/images", &buf)
	req.Header.Set("Content-Type", w2.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadImage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUploadHandler_ServeAsset_Success(t *testing.T) {
	dir := t.TempDir()
	filename := "default_image.jpg"
	expected := []byte("FAKEJPEG")
	if err := os.WriteFile(filepath.Join(dir, filename), expected, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	h := handlers.NewUploadHandler(dir)
	req := chiRequest(http.MethodGet, "/api/assets/"+filename, nil, map[string]string{"filename": filename})
	w := httptest.NewRecorder()

	h.ServeAsset(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	got, _ := io.ReadAll(w.Body)
	if string(got) != string(expected) {
		t.Errorf("body mismatch")
	}
	if cc := w.Header().Get("Cache-Control"); cc == "" {
		t.Errorf("expected Cache-Control header")
	}
}

func TestUploadHandler_ServeAsset_TraversalRejected(t *testing.T) {
	dir := t.TempDir()
	h := handlers.NewUploadHandler(dir)

	req := chiRequest(http.MethodGet, "/api/assets/..%2Fsecret", nil, map[string]string{"filename": "../secret"})
	w := httptest.NewRecorder()

	h.ServeAsset(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUploadHandler_ServeAsset_NotFound(t *testing.T) {
	dir := t.TempDir()
	h := handlers.NewUploadHandler(dir)

	req := chiRequest(http.MethodGet, "/api/assets/missing.jpg", nil, map[string]string{"filename": "missing.jpg"})
	w := httptest.NewRecorder()

	h.ServeAsset(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestEnsureDefaultImage_CopiesIfMissing(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()
	src := filepath.Join(srcDir, "default_image.jpg")
	if err := os.WriteFile(src, []byte("seed"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := handlers.EnsureDefaultImage(dstDir, src); err != nil {
		t.Fatalf("EnsureDefaultImage: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dstDir, "default_image.jpg"))
	if err != nil {
		t.Fatalf("read dest: %v", err)
	}
	if string(got) != "seed" {
		t.Errorf("unexpected content: %s", got)
	}
}

func TestEnsureDefaultImage_NoOpIfPresent(t *testing.T) {
	dstDir := t.TempDir()
	dest := filepath.Join(dstDir, "default_image.jpg")
	if err := os.WriteFile(dest, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := handlers.EnsureDefaultImage(dstDir, "/does/not/exist.jpg"); err != nil {
		t.Fatalf("EnsureDefaultImage should be a no-op when destination exists: %v", err)
	}

	got, _ := os.ReadFile(dest)
	if string(got) != "existing" {
		t.Errorf("file was overwritten")
	}
}
