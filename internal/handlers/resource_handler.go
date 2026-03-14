package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hema-lessons/internal/pagination"
	"hema-lessons/internal/store"
)

type ResourceHandler struct {
	store *store.Store
}

func NewResourceHandler(s *store.Store) *ResourceHandler {
	return &ResourceHandler{store: s}
}

func (h *ResourceHandler) List(w http.ResponseWriter, r *http.Request) {
	params := pagination.ParseParams(r)

	resources, totalCount := h.store.ListResources(params)

	response := pagination.NewResponse(resources, params, totalCount)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ResourceHandler) Get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	prefix := "/api/resources/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	idStr := path[len(prefix):]
	if idStr == "" {
		http.Error(w, "invalid resource ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid resource ID", http.StatusBadRequest)
		return
	}

	resource := h.store.GetResourceByID(id)
	if resource == nil {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resource); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
