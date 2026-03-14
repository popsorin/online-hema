package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hema-lessons/internal/store"
)

type ItemHandler struct {
	store *store.Store
}

func NewItemHandler(s *store.Store) *ItemHandler {
	return &ItemHandler{store: s}
}

// ListBySection handles GET /api/sections/:id/items — returns items within a section.
func (h *ItemHandler) ListBySection(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	prefix := "/api/sections/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	remaining := path[len(prefix):]
	parts := strings.Split(remaining, "/")

	if len(parts) != 2 || parts[1] != "items" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	sectionID, err := strconv.Atoi(parts[0])
	if err != nil || sectionID <= 0 {
		http.Error(w, "invalid section ID", http.StatusBadRequest)
		return
	}

	if h.store.GetSectionByID(sectionID) == nil {
		http.Error(w, "section not found", http.StatusNotFound)
		return
	}

	items := h.store.ListItemsBySectionID(sectionID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

