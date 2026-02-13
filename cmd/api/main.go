package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"hema-lessons/internal/config"
	"hema-lessons/internal/database"
	"hema-lessons/internal/handlers"
	"hema-lessons/internal/repository"
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

	db, err := database.ConnectWithDSN(cfg.GetDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("database connection established")

	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("database migrations completed")

	if err := seedSampleDataIfEmpty(db); err != nil {
		log.Printf("warning: failed to seed sample data: %v", err)
	} else {
		log.Println("database sample data seeded")
	}

	fightingBookRepo := repository.NewFightingBookRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	techniqueRepo := repository.NewTechniqueRepository(db)

	fightingBookHandler := handlers.NewFightingBookHandler(fightingBookRepo)
	chapterHandler := handlers.NewChapterHandler(chapterRepo)
	techniqueHandler := handlers.NewTechniqueHandler(techniqueRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/healthz" {
			healthzHandler(db, cfg)(w, r)
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

func seedSampleDataIfEmpty(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sword_masters").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = db.Exec(`
		INSERT INTO sword_masters (name, bio, birth_year, death_year) VALUES
		('Johannes Liechtenauer', 'German fencing master, founder of the Liechtenauer tradition of German fencing', 1300, 1389),
		('Fiore dei Liberi', 'Italian fencing master who wrote one of the most important medieval fighting treatises', 1350, 1420),
		('Sigmund Ringeck', 'German fencing master and student of Johannes Liechtenauer', 1400, 1470)
	`)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT id, name FROM sword_masters ORDER BY name")
	if err != nil {
		return err
	}
	defer rows.Close()

	var masterIDs []int
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		masterIDs = append(masterIDs, id)
	}

	_, err = db.Exec(`
		INSERT INTO fighting_books (sword_master_id, title, description, publication_year) VALUES
		($1, 'Zettel', 'A mnemonic poem describing the core principles of Liechtenauer''s fighting system', 1389),
		($2, 'Fior di Battaglia', 'The Flower of Battle - a comprehensive medieval combat manual covering armed and unarmed combat', 1409),
		($3, 'Fechtbuch', 'A detailed commentary on Liechtenauer''s teachings with practical applications', 1440)
	`, masterIDs[0], masterIDs[1], masterIDs[2])
	if err != nil {
		return err
	}

	rows, err = db.Query(`
		SELECT fb.id, fb.title, sm.id as master_id
		FROM fighting_books fb
		JOIN sword_masters sm ON fb.sword_master_id = sm.id
		ORDER BY fb.title
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bookData := make(map[string]int)
	for rows.Next() {
		var id int
		var title string
		var masterID int
		if err := rows.Scan(&id, &title, &masterID); err != nil {
			return err
		}
		bookData[title] = id
	}

	fiorID := bookData["Fior di Battaglia"]
	_, err = db.Exec(`
		INSERT INTO chapters (fighting_book_id, chapter_number, title, description) VALUES
		($1, 1, 'Wrestling', 'Techniques for unarmed combat and grappling'),
		($1, 2, 'Dagger Combat', 'Fighting with the dagger in various situations'),
		($1, 3, 'Longsword', 'The art of fighting with the longsword'),
		($1, 4, 'Poleaxe', 'Combat techniques with the poleaxe')
	`, fiorID)
	if err != nil {
		return err
	}

	var longswordChapterID int
	err = db.QueryRow(`
		SELECT id FROM chapters
		WHERE fighting_book_id = $1 AND chapter_number = 3
	`, fiorID).Scan(&longswordChapterID)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO techniques (chapter_id, name, description, instructions, order_in_chapter) VALUES
		($1, 'First Guard - Posta di Donna', 'The Woman''s Guard - a high guard position', 'Hold the sword with the hilt near your right shoulder, point aimed at the opponent''s face', 1),
		($1, 'Zornhau', 'The Wrath Strike - a powerful diagonal cut', 'Strike diagonally from your right shoulder to the opponent''s left side with full commitment', 2),
		($1, 'Krumphau', 'The Crooked Strike - an off-line attack', 'Step offline and strike with the false edge, hands crossed', 3)
	`, longswordChapterID)

	return err
}

func healthzHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "database unhealthy", http.StatusServiceUnavailable)
			return
		}

		resp := healthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		if cfg.IsDevelopment() {
			resp.Status = fmt.Sprintf("%s (%s)", resp.Status, cfg.App.Environment)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
