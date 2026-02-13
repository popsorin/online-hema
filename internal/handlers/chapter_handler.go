package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hema-lessons/internal/repository"
)

type ChapterHandler struct {
	repo *repository.ChapterRepository
}

func NewChapterHandler(repo *repository.ChapterRepository) *ChapterHandler {
	return &ChapterHandler{repo: repo}
}

func (h *ChapterHandler) ListByFightingBook(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	prefix := "/api/fighting-books/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	remaining := path[len(prefix):]
	parts := strings.Split(remaining, "/")

	if len(parts) != 2 || parts[1] != "chapters" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	idStr := parts[0]
	if idStr == "" {
		http.Error(w, "invalid fighting book ID", http.StatusBadRequest)
		return
	}

	fightingBookID, err := strconv.Atoi(idStr)
	if err != nil || fightingBookID <= 0 {
		http.Error(w, "invalid fighting book ID", http.StatusBadRequest)
		return
	}

	exists, err := h.repo.FightingBookExists(fightingBookID)
	if err != nil {
		log.Printf("failed to check fighting book existence: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "fighting book not found", http.StatusNotFound)
		return
	}

	chapters, err := h.repo.ListByFightingBookID(fightingBookID)
	if err != nil {
		log.Printf("failed to list chapters: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chapters); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
