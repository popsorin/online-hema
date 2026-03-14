package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hema-lessons/internal/store"
)

type SectionHandler struct {
	store *store.Store
}

func NewSectionHandler(s *store.Store) *SectionHandler {
	return &SectionHandler{store: s}
}

// ListByResource handles GET /api/resources/:id/sections — returns root-level sections for a resource.
func (h *SectionHandler) ListByBook(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	prefix := "/api/resources/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	remaining := path[len(prefix):]
	parts := strings.Split(remaining, "/")

	if len(parts) != 2 || parts[1] != "sections" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	resourceID, err := strconv.Atoi(parts[0])
	if err != nil || resourceID <= 0 {
		http.Error(w, "invalid resource ID", http.StatusBadRequest)
		return
	}

	if !h.store.ResourceExists(resourceID) {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	sections := h.store.ListRootSectionsByResourceID(resourceID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sections); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// Get handles GET /api/sections/:id — returns a single section.
func (h *SectionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseSectionID(w, r.URL.Path)
	if !ok {
		return
	}

	section := h.store.GetSectionByID(id)
	if section == nil {
		http.Error(w, "section not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(section); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// ListChildren handles GET /api/sections/:id/sections — returns child sections.
func (h *SectionHandler) ListChildren(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	prefix := "/api/sections/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	remaining := path[len(prefix):]
	parts := strings.Split(remaining, "/")

	if len(parts) != 2 || parts[1] != "sections" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	parentID, err := strconv.Atoi(parts[0])
	if err != nil || parentID <= 0 {
		http.Error(w, "invalid section ID", http.StatusBadRequest)
		return
	}

	if h.store.GetSectionByID(parentID) == nil {
		http.Error(w, "section not found", http.StatusNotFound)
		return
	}

	sections := h.store.ListChildSections(parentID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sections); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// parseSectionID extracts an integer ID from a path of the form /api/sections/:id
// (without a trailing segment). Returns false and writes an error response if invalid.
func parseSectionID(w http.ResponseWriter, path string) (int, bool) {
	prefix := "/api/sections/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return 0, false
	}

	idStr := path[len(prefix):]
	if idStr == "" {
		http.Error(w, "invalid section ID", http.StatusBadRequest)
		return 0, false
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid section ID", http.StatusBadRequest)
		return 0, false
	}

	return id, true
}
