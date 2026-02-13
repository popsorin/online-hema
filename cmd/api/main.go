package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"hema-lessons/internal/config"
	"hema-lessons/internal/handlers"
	"hema-lessons/internal/store"
)

type healthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	if cfg.IsDevelopment() {
		log.Printf("starting in development mode")
	}

	dataStore, err := store.New()
	if err != nil {
		log.Fatalf("failed to load data store: %v", err)
	}
	log.Println("data store loaded")

	fightingBookHandler := handlers.NewFightingBookHandler(dataStore)
	chapterHandler := handlers.NewChapterHandler(dataStore)
	techniqueHandler := handlers.NewTechniqueHandler(dataStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/healthz" {
			healthzHandler(cfg)(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/fighting-books") {
			if path == "/api/fighting-books" || path == "/api/fighting-books/" {
				if r.Method == http.MethodGet {
					fightingBookHandler.List(w, r)
					return
				}
			} else if strings.HasPrefix(path, "/api/fighting-books/") && path != "/api/fighting-books/" {
				if strings.HasSuffix(path, "/chapters") {
					if r.Method == http.MethodGet {
						chapterHandler.ListByFightingBook(w, r)
						return
					}
				} else if r.Method == http.MethodGet {
					fightingBookHandler.Get(w, r)
					return
				}
			}
		}

		if strings.HasPrefix(path, "/api/chapters/") {
			if strings.HasSuffix(path, "/techniques") {
				if r.Method == http.MethodGet {
					techniqueHandler.ListByChapter(w, r)
					return
				}
			}
		}

		http.NotFound(w, r)
	})

	server := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(cfg.Server.ReadHeaderTimeout) * time.Second,
	}

	log.Printf("hema api listening on %s", cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func healthzHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := healthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		if cfg.IsDevelopment() {
			resp.Status = resp.Status + " (" + cfg.App.Environment + ")"
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
