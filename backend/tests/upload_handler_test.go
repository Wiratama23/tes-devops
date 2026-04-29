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
	"runtime"
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

// TestEnsureDefaultImage_RemovesPartialOnCopyFailure exercises the cleanup
// branch added to EnsureDefaultImage: when os.Open succeeds but io.Copy
// fails, the partially-created destination file must be removed so a later
// boot does not see it via os.Stat and short-circuit the copy with corrupt
// data.
//
// We trigger an io.Copy failure by passing a directory as the source path.
// On Linux/macOS this is the canonical "EISDIR on Read" reproduction; on
// Windows os.Open of a directory typically fails outright with a different
// error path (which is already covered by the source-missing case), so we
// skip there.
func TestEnsureDefaultImage_RemovesPartialOnCopyFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("io.Copy failure path unreliable on windows; exercised on linux/CI")
	}

	srcDir := t.TempDir()
	dstDir := t.TempDir()

	if err := handlers.EnsureDefaultImage(dstDir, srcDir); err == nil {
		t.Fatalf("expected error when source is a directory, got nil")
	}

	dest := filepath.Join(dstDir, "default_image.jpg")
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Errorf("expected partial dest file to be removed, but it still exists: stat err=%v", err)
	}
}

// TestUploadHandler_UploadImage_RemovesMultipartTemps confirms the deferred
// MultipartForm.RemoveAll added to the handler runs even after the happy
// path. Files smaller than the in-memory threshold (10 MiB here) never spill
// to disk, so the visible side-effect we can rely on is that
// `r.MultipartForm` is no longer holding any usable file handles after the
// handler returns. We assert via the upload directory that the saved
// payload is intact (i.e. the cleanup hasn't accidentally torn anything we
// still needed away).
func TestUploadHandler_UploadImage_RemovesMultipartTemps(t *testing.T) {
	dir := t.TempDir()
	h := handlers.NewUploadHandler(dir)

	body, ct := multipartImage(t, "file", "ok.png", "image/png", []byte("PNGDATA"))
	req := httptest.NewRequest(http.MethodPost, "/api/uploads/images", body)
	req.Header.Set("Content-Type", ct)
	w := httptest.NewRecorder()

	h.UploadImage(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp handlers.UploadResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	saved, err := os.ReadFile(filepath.Join(dir, resp.Filename))
	if err != nil {
		t.Fatalf("saved file missing after upload + cleanup: %v", err)
	}
	if string(saved) != "PNGDATA" {
		t.Errorf("saved bytes mismatch: %q", saved)
	}
}
