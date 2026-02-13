package testutil

import (
	"database/sql"
	"log"
	"testing"
)

func SeedSwordMasters(t *testing.T, db *sql.DB) []int {
	t.Helper()

	masters := []struct {
		name      string
		bio       string
		birthYear int
		deathYear int
	}{
		{"Test Master 1", "First test master", 1300, 1380},
		{"Test Master 2", "Second test master", 1350, 1420},
		{"Test Master 3", "Third test master", 1400, 1470},
	}

	var ids []int
	for _, m := range masters {
		var id int
		err := db.QueryRow(`
			INSERT INTO sword_masters (name, bio, birth_year, death_year)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, m.name, m.bio, m.birthYear, m.deathYear).Scan(&id)

		if err != nil {
			t.Fatalf("failed to seed sword master: %v", err)
		}
		ids = append(ids, id)
	}

	return ids
}

func SeedFightingBooks(t *testing.T, db *sql.DB, masterIDs []int) []int {
	t.Helper()

	books := []struct {
		masterID    int
		title       string
		description string
		pubYear     int
	}{
		{masterIDs[0], "Book A", "First test book", 1400},
		{masterIDs[0], "Book B", "Second test book", 1410},
		{masterIDs[1], "Book C", "Third test book", 1420},
		{masterIDs[1], "Book D", "Fourth test book", 1430},
		{masterIDs[2], "Book E", "Fifth test book", 1440},
	}

	var ids []int
	for _, b := range books {
		var id int
		err := db.QueryRow(`
			INSERT INTO fighting_books (sword_master_id, title, description, publication_year)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, b.masterID, b.title, b.description, b.pubYear).Scan(&id)

		if err != nil {
			t.Fatalf("failed to seed fighting book: %v", err)
		}
		ids = append(ids, id)
	}

	return ids
}

func SeedSampleData(t *testing.T, db *sql.DB) {
	t.Helper()
	seedSampleData(db, t.Fatalf)
}

func SeedSampleDataForApp(db *sql.DB) error {
	return seedSampleData(db, func(format string, args ...interface{}) {
		log.Printf("Error seeding data: "+format, args...)
	})
}

func seedSampleData(db *sql.DB, failFunc func(string, ...interface{})) error {
	_, err := db.Exec(`
		INSERT INTO sword_masters (name, bio, birth_year, death_year) VALUES
		('Johannes Liechtenauer', 'German fencing master, founder of the Liechtenauer tradition of German fencing', 1300, 1389),
		('Fiore dei Liberi', 'Italian fencing master who wrote one of the most important medieval fighting treatises', 1350, 1420),
		('Sigmund Ringeck', 'German fencing master and student of Johannes Liechtenauer', 1400, 1470)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		failFunc("failed to seed sample sword masters: %v", err)
		return err
	}

	rows, err := db.Query("SELECT id, name FROM sword_masters WHERE name IN ('Johannes Liechtenauer', 'Fiore dei Liberi', 'Sigmund Ringeck') ORDER BY name")
	if err != nil {
		failFunc("failed to get sword master IDs: %v", err)
		return err
	}
	defer rows.Close()

	var masterIDs []int
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			failFunc("failed to scan sword master: %v", err)
			return err
		}
		masterIDs = append(masterIDs, id)
	}

	if len(masterIDs) != 3 {
		failFunc("expected 3 sword masters, got %d", len(masterIDs))
		return err
	}

	_, err = db.Exec(`
		INSERT INTO fighting_books (sword_master_id, title, description, publication_year) VALUES
		($1, 'Zettel', 'A mnemonic poem describing the core principles of Liechtenauer''s fighting system', 1389),
		($2, 'Fior di Battaglia', 'The Flower of Battle - a comprehensive medieval combat manual covering armed and unarmed combat', 1409),
		($3, 'Fechtbuch', 'A detailed commentary on Liechtenauer''s teachings with practical applications', 1440)
		ON CONFLICT (title) DO NOTHING
	`, masterIDs[0], masterIDs[1], masterIDs[2])
	if err != nil {
		failFunc("failed to seed sample fighting books: %v", err)
		return err
	}

	log.Printf("Successfully seeded sample data: %d sword masters and 3 fighting books", len(masterIDs))
	return nil
}

func SeedChapters(t *testing.T, db *sql.DB, bookIDs []int) []int {
	t.Helper()

	chapters := []struct {
		bookID        int
		chapterNumber int
		title         string
		description   string
	}{
		{bookIDs[0], 1, "Chapter 1", "First chapter of Book A"},
		{bookIDs[0], 2, "Chapter 2", "Second chapter of Book A"},
		{bookIDs[0], 3, "Chapter 3", "Third chapter of Book A"},
		{bookIDs[1], 1, "Introduction", "First chapter of Book B"},
		{bookIDs[1], 2, "Advanced Techniques", "Second chapter of Book B"},
	}

	var ids []int
	for _, c := range chapters {
		var id int
		err := db.QueryRow(`
			INSERT INTO chapters (fighting_book_id, chapter_number, title, description)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, c.bookID, c.chapterNumber, c.title, c.description).Scan(&id)

		if err != nil {
			t.Fatalf("failed to seed chapter: %v", err)
		}
		ids = append(ids, id)
	}

	return ids
}

func SeedTechniques(t *testing.T, db *sql.DB, chapterIDs []int) []int {
	t.Helper()

	techniques := []struct {
		chapterID      int
		name           string
		description    string
		instructions   string
		orderInChapter int
	}{
		{chapterIDs[0], "Technique 1", "First technique", "Step 1, Step 2", 1},
		{chapterIDs[0], "Technique 2", "Second technique", "Step A, Step B", 2},
		{chapterIDs[0], "Technique 3", "Third technique", "Instructions here", 3},
		{chapterIDs[1], "Basic Move", "A basic move", "Do this", 1},
		{chapterIDs[1], "Advanced Move", "An advanced move", "Then do this", 2},
	}

	var ids []int
	for _, tech := range techniques {
		var id int
		err := db.QueryRow(`
			INSERT INTO techniques (chapter_id, name, description, instructions, order_in_chapter)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, tech.chapterID, tech.name, tech.description, tech.instructions, tech.orderInChapter).Scan(&id)

		if err != nil {
			t.Fatalf("failed to seed technique: %v", err)
		}
		ids = append(ids, id)
	}

	return ids
}
