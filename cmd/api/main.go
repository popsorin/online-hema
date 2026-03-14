package main

import (
	"encoding/json"
	"fmt"
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

// #region agent log
func debugLog(location, message, hypothesisID string, data map[string]interface{}) {
	f, err := os.OpenFile("/var/www/hema-lessons/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	d, _ := json.Marshal(data)
	fmt.Fprintf(f, `{"timestamp":%d,"location":%q,"message":%q,"hypothesisId":%q,"data":%s}`+"\n",
		time.Now().UnixMilli(), location, message, hypothesisID, string(d))
}

// #endregion

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
		// #region agent log
		debugLog("main.go:store.New", "store load FAILED", "H-B", map[string]interface{}{"error": err.Error()})
		// #endregion
		slog.Error("failed to load data store", "error", err)
		os.Exit(1)
	}
	// #region agent log
	debugLog("main.go:store.New", "store loaded OK - new binary is running", "H-A", map[string]interface{}{"built": "post-rename"})
	// #endregion
	slog.Info("data store loaded")

	resourceHandler := handlers.NewResourceHandler(dataStore)
	sectionHandler := handlers.NewSectionHandler(dataStore)
	itemHandler := handlers.NewItemHandler(dataStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// #region agent log
		debugLog("main.go:router", "incoming request", "H-D", map[string]interface{}{"method": r.Method, "path": path})
		// #endregion

		if strings.HasPrefix(path, "/assets/") {
			http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))).ServeHTTP(w, r)
			return
		}

		if path == "/healthz" {
			healthzHandler(cfg)(w, r)
			return
		}

		// GET /api/resources
		if path == "/api/resources" || path == "/api/resources/" {
			if r.Method == http.MethodGet {
				resourceHandler.List(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/api/resources/") {
			remaining := path[len("/api/resources/"):]
			parts := strings.SplitN(remaining, "/", 2)

			if len(parts) == 1 {
				// GET /api/resources/:id
				if r.Method == http.MethodGet {
					resourceHandler.Get(w, r)
					return
				}
			} else if len(parts) == 2 && parts[1] == "sections" {
				// GET /api/resources/:id/sections
				if r.Method == http.MethodGet {
					sectionHandler.ListByBook(w, r)
					return
				}
			}
		}

		if strings.HasPrefix(path, "/api/sections/") {
			remaining := path[len("/api/sections/"):]
			parts := strings.SplitN(remaining, "/", 2)

			if len(parts) == 1 {
				// GET /api/sections/:id
				if r.Method == http.MethodGet {
					sectionHandler.Get(w, r)
					return
				}
			} else if len(parts) == 2 {
				switch parts[1] {
				case "sections":
					// GET /api/sections/:id/sections
					if r.Method == http.MethodGet {
						sectionHandler.ListChildren(w, r)
						return
					}
				case "items":
					// GET /api/sections/:id/items
					if r.Method == http.MethodGet {
						itemHandler.ListBySection(w, r)
						return
					}
				}
			}
		}

		// #region agent log
		debugLog("main.go:notfound", "request fell through to 404", "H-D", map[string]interface{}{"path": path, "method": r.Method})
		// #endregion
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
