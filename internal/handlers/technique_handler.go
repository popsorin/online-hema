package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hema-lessons/internal/store"
)

type TechniqueHandler struct {
	store *store.Store
}

func NewTechniqueHandler(s *store.Store) *TechniqueHandler {
	return &TechniqueHandler{store: s}
}

func (h *TechniqueHandler) ListByChapter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	prefix := "/api/chapters/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	remaining := path[len(prefix):]
	parts := strings.Split(remaining, "/")

	if len(parts) != 2 || parts[1] != "techniques" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	idStr := parts[0]
	if idStr == "" {
		http.Error(w, "invalid chapter ID", http.StatusBadRequest)
		return
	}

	chapterID, err := strconv.Atoi(idStr)
	if err != nil || chapterID <= 0 {
		http.Error(w, "invalid chapter ID", http.StatusBadRequest)
		return
	}

	if !h.store.ChapterExists(chapterID) {
		http.Error(w, "chapter not found", http.StatusNotFound)
		return
	}

	techniques := h.store.ListTechniquesByChapterID(chapterID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(techniques); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
