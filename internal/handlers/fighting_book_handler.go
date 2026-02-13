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

type FightingBookHandler struct {
	store *store.Store
}

func NewFightingBookHandler(s *store.Store) *FightingBookHandler {
	return &FightingBookHandler{store: s}
}

func (h *FightingBookHandler) List(w http.ResponseWriter, r *http.Request) {
	params := pagination.ParseParams(r)

	books, totalCount := h.store.ListFightingBooks(params)

	response := pagination.NewResponse(books, params, totalCount)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *FightingBookHandler) Get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	prefix := "/api/fighting-books/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	if len(path) <= len(prefix) {
		http.Error(w, "invalid fighting book ID", http.StatusBadRequest)
		return
	}

	idStr := path[len(prefix):]
	if idStr == "" {
		http.Error(w, "invalid fighting book ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid fighting book ID", http.StatusBadRequest)
		return
	}

	book := h.store.GetFightingBookByID(id)
	if book == nil {
		http.Error(w, "fighting book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
