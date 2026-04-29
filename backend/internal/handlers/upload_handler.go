package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// UploadHandler stores admin-uploaded images in UPLOADS_DIR and serves them
// back through GET /api/assets/{filename}. The handler intentionally avoids
// any database tracking — the canonical record of an image is the
// `image_path` column on the product row.
type UploadHandler struct {
	uploadsDir   string
	maxBytes     int64
	allowedExt   map[string]bool
	allowedMimes map[string]bool
}

func NewUploadHandler(uploadsDir string) *UploadHandler {
	return &UploadHandler{
		uploadsDir: uploadsDir,
		maxBytes:   10 << 20, // 10 MiB
		allowedExt: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
			".gif":  true,
		},
		allowedMimes: map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
			"image/gif":  true,
		},
	}
}

type UploadResponse struct {
	Path     string `json:"path"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// UploadImage handles POST /api/uploads/images (admin-only, multipart).
// The form field is `file`. Filenames are randomized to prevent collisions
// and path traversal.
func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxBytes)
	if err := r.ParseMultipartForm(h.maxBytes); err != nil {
		http.Error(w, "file too large or invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !h.allowedExt[ext] {
		http.Error(w, "unsupported file extension", http.StatusBadRequest)
		return
	}

	if ct := header.Header.Get("Content-Type"); ct != "" && !h.allowedMimes[strings.ToLower(ct)] {
		http.Error(w, "unsupported content type", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll(h.uploadsDir, 0o755); err != nil {
		http.Error(w, "failed to prepare uploads directory", http.StatusInternalServerError)
		return
	}

	filename, err := randomFilename(ext)
	if err != nil {
		http.Error(w, "failed to generate filename", http.StatusInternalServerError)
		return
	}

	dstPath := filepath.Join(h.uploadsDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		_ = os.Remove(dstPath)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	resp := UploadResponse{
		Path:     "assets/" + filename,
		URL:      "/api/assets/" + filename,
		Filename: filename,
		Size:     written,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ServeAsset handles GET /api/assets/{filename}. It enforces no path traversal
// and sets a reasonable Cache-Control header so nginx + the browser cache it.
func (h *UploadHandler) ServeAsset(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "filename")
	if name == "" {
		http.NotFound(w, r)
		return
	}

	cleaned := filepath.Base(name)
	if cleaned != name || strings.Contains(name, "..") {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	full := filepath.Join(h.uploadsDir, cleaned)
	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}

	ext := strings.ToLower(filepath.Ext(cleaned))
	if ct := mime.TypeByExtension(ext); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeFile(w, r, full)
}

// EnsureDefaultImage copies the bundled default image into the uploads
// directory if it is missing. Called from main.go on boot so that
// `GET /api/assets/default_image.jpg` always resolves, even before the
// admin uploads anything.
func EnsureDefaultImage(uploadsDir, sourcePath string) error {
	if uploadsDir == "" {
		return fmt.Errorf("uploadsDir is empty")
	}
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		return fmt.Errorf("create uploads dir: %w", err)
	}

	dest := filepath.Join(uploadsDir, "default_image.jpg")
	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	src, err := os.Open(sourcePath)
	if err != nil {
		// Source missing is not fatal — uploads dir is created, just no default.
		return fmt.Errorf("open default image source: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create default image: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy default image: %w", err)
	}
	return nil
}

func randomFilename(ext string) (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	stamp := time.Now().UTC().Format("20060102")
	return fmt.Sprintf("%s-%s%s", stamp, hex.EncodeToString(buf), ext), nil
}
