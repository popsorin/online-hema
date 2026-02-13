package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hema-lessons/internal/pagination"
	"hema-lessons/internal/repository"
)

type FightingBookHandler struct {
	repo *repository.FightingBookRepository
}

func NewFightingBookHandler(repo *repository.FightingBookRepository) *FightingBookHandler {
	return &FightingBookHandler{repo: repo}
}

func (h *FightingBookHandler) List(w http.ResponseWriter, r *http.Request) {
	params := pagination.ParseParams(r)

	books, totalCount, err := h.repo.List(params)
	if err != nil {
		log.Printf("failed to list fighting books: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

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

	book, err := h.repo.GetByID(id)
	if err != nil {
		log.Printf("failed to get fighting book: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

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
