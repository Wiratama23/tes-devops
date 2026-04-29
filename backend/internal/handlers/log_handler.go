package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// LogHandler accepts client-side error reports forwarded by the frontend's
// central logger. It is intentionally minimal: payloads are validated lightly
// and emitted to the standard server log so they are visible alongside other
// API logs (and can be picked up by any log aggregator the WAF/nginx already
// forward).
type LogHandler struct {
	maxMessage int
	maxStack   int
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		maxMessage: 2000,
		maxStack:   8000,
	}
}

type ClientLogEntry struct {
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Stack     string         `json:"stack,omitempty"`
	URL       string         `json:"url,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`
}

// Receive handles POST /api/logs.
func (h *LogHandler) Receive(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB cap
	var entry ClientLogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	level := strings.ToLower(strings.TrimSpace(entry.Level))
	switch level {
	case "info", "warn", "error":
	default:
		level = "error"
	}

	if entry.UserAgent == "" {
		entry.UserAgent = r.Header.Get("User-Agent")
	}

	msg := truncate(entry.Message, h.maxMessage)
	stack := truncate(entry.Stack, h.maxStack)

	meta, _ := json.Marshal(entry.Meta)

	log.Printf(
		"[client-%s] %s url=%q ua=%q meta=%s stack=%q",
		level,
		msg,
		entry.URL,
		entry.UserAgent,
		string(meta),
		stack,
	)

	w.WriteHeader(http.StatusNoContent)
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "...<truncated>"
}
