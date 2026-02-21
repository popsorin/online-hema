package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"hema-lessons/internal/config"
	"hema-lessons/internal/handlers"
	"hema-lessons/internal/middleware"
	"hema-lessons/internal/store"
)

type healthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Set up structured logging - JSON in production, text in development
	var logHandler slog.Handler
	if cfg.IsProduction() {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	slog.SetDefault(slog.New(logHandler))

	slog.Info("starting application", "environment", cfg.App.Environment)

	// Initialize Sentry for error tracking
	if cfg.App.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.App.SentryDSN,
			Environment:      cfg.App.Environment,
			TracesSampleRate: 0.1, // 10% of transactions for performance monitoring
		})
		if err != nil {
			slog.Error("failed to initialize Sentry", "error", err)
		} else {
			slog.Info("sentry initialized")
			defer sentry.Flush(2 * time.Second)
		}
	}

	dataStore, err := store.New()
	if err != nil {
		slog.Error("failed to load data store", "error", err)
		os.Exit(1)
	}
	slog.Info("data store loaded")

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

	// Wrap handler with middleware (order: Recovery -> RequestLogger -> mux)
	httpHandler := middleware.Recovery(middleware.RequestLogger(mux))

	server := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           httpHandler,
		ReadHeaderTimeout: time.Duration(cfg.Server.ReadHeaderTimeout) * time.Second,
	}

	slog.Info("hema api listening", "addr", cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
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
