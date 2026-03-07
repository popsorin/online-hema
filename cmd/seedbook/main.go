package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"hema-lessons/internal/models"
)

const dataDir = "internal/store/data"

// --- YAML input structures ---

type BookDefinition struct {
	SwordMaster    string           `yaml:"sword_master"`
	NewSwordMaster *NewSwordMaster  `yaml:"new_sword_master,omitempty"`
	Book           BookInput        `yaml:"book"`
	Chapters       []ChapterInput   `yaml:"chapters"`
}

type NewSwordMaster struct {
	Name      string `yaml:"name"`
	Bio       string `yaml:"bio"`
	BirthYear *int   `yaml:"birth_year,omitempty"`
	DeathYear *int   `yaml:"death_year,omitempty"`
}

type BookInput struct {
	Title           string `yaml:"title"`
	Description     string `yaml:"description"`
	PublicationYear *int   `yaml:"publication_year,omitempty"`
}

type ChapterInput struct {
	Title       string           `yaml:"title"`
	Description string           `yaml:"description"`
	Techniques  []TechniqueInput `yaml:"techniques"`
}

type TechniqueInput struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Instructions string `yaml:"instructions"`
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Print what would be added without writing files")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: seedbook [--dry-run] <path-to-yaml>\n")
		os.Exit(1)
	}
	yamlPath := flag.Arg(0)

	raw, err := os.ReadFile(yamlPath)
	if err != nil {
		fatal("reading YAML file: %v", err)
	}

	var def BookDefinition
	if err := yaml.Unmarshal(raw, &def); err != nil {
		fatal("parsing YAML: %v", err)
	}

	if err := validate(&def); err != nil {
		fatal("validation: %v", err)
	}

	masters, err := loadJSON[models.SwordMaster](filepath.Join(dataDir, "sword_masters.json"))
	if err != nil {
		fatal("%v", err)
	}
	books, err := loadJSON[models.FightingBook](filepath.Join(dataDir, "fighting_books.json"))
	if err != nil {
		fatal("%v", err)
	}
	chapters, err := loadJSON[models.Chapter](filepath.Join(dataDir, "chapters.json"))
	if err != nil {
		fatal("%v", err)
	}
	techniques, err := loadJSON[models.Technique](filepath.Join(dataDir, "techniques.json"))
	if err != nil {
		fatal("%v", err)
	}

	masterID, newMaster := resolveSwordMaster(&def, masters)
	if masterID == 0 && newMaster == nil {
		fatal("sword master %q not found and no new_sword_master block provided", def.SwordMaster)
	}

	if dupBook := findBookByTitle(books, def.Book.Title); dupBook != nil {
		fatal("book %q already exists (id=%d)", def.Book.Title, dupBook.ID)
	}

	// Build new entities with auto-incremented IDs
	nextMasterID := maxID(masters, func(m models.SwordMaster) int { return m.ID }) + 1
	nextBookID := maxID(books, func(b models.FightingBook) int { return b.ID }) + 1
	nextChapterID := maxID(chapters, func(c models.Chapter) int { return c.ID }) + 1
	nextTechniqueID := maxID(techniques, func(t models.Technique) int { return t.ID }) + 1

	if newMaster != nil {
		newMaster.ID = nextMasterID
		masterID = nextMasterID
		masters = append(masters, *newMaster)
		fmt.Printf("[+] Sword Master: id=%d name=%q\n", newMaster.ID, newMaster.Name)
	}

	newBook := models.FightingBook{
		ID:              nextBookID,
		SwordMasterID:   masterID,
		Title:           def.Book.Title,
		Description:     def.Book.Description,
		PublicationYear: def.Book.PublicationYear,
	}
	books = append(books, newBook)
	fmt.Printf("[+] Fighting Book: id=%d title=%q (master_id=%d)\n", newBook.ID, newBook.Title, masterID)

	for i, ch := range def.Chapters {
		newChapter := models.Chapter{
			ID:             nextChapterID,
			FightingBookID: nextBookID,
			ChapterNumber:  i + 1,
			Title:          ch.Title,
			Description:    ch.Description,
		}
		chapters = append(chapters, newChapter)
		fmt.Printf("[+] Chapter: id=%d #%d %q (%d techniques)\n", newChapter.ID, newChapter.ChapterNumber, newChapter.Title, len(ch.Techniques))

		for j, tech := range ch.Techniques {
			newTech := models.Technique{
				ID:             nextTechniqueID,
				ChapterID:      nextChapterID,
				Name:           tech.Name,
				Description:    tech.Description,
				Instructions:   tech.Instructions,
				OrderInChapter: j + 1,
			}
			techniques = append(techniques, newTech)
			fmt.Printf("    [+] Technique: id=%d #%d %q\n", newTech.ID, newTech.OrderInChapter, newTech.Name)
			nextTechniqueID++
		}
		nextChapterID++
	}

	if *dryRun {
		fmt.Println("\n--dry-run: no files written.")
		return
	}

	if err := writeJSON(filepath.Join(dataDir, "sword_masters.json"), masters); err != nil {
		fatal("writing sword_masters.json: %v", err)
	}
	if err := writeJSON(filepath.Join(dataDir, "fighting_books.json"), books); err != nil {
		fatal("writing fighting_books.json: %v", err)
	}
	if err := writeJSON(filepath.Join(dataDir, "chapters.json"), chapters); err != nil {
		fatal("writing chapters.json: %v", err)
	}
	if err := writeJSON(filepath.Join(dataDir, "techniques.json"), techniques); err != nil {
		fatal("writing techniques.json: %v", err)
	}

	fmt.Println("\nAll data files updated successfully.")
}

func validate(def *BookDefinition) error {
	if def.SwordMaster == "" && def.NewSwordMaster == nil {
		return fmt.Errorf("either sword_master or new_sword_master is required")
	}
	if def.Book.Title == "" {
		return fmt.Errorf("book.title is required")
	}
	if len(def.Chapters) == 0 {
		return fmt.Errorf("at least one chapter is required")
	}
	for i, ch := range def.Chapters {
		if ch.Title == "" {
			return fmt.Errorf("chapter[%d].title is required", i)
		}
		for j, t := range ch.Techniques {
			if t.Name == "" {
				return fmt.Errorf("chapter[%d].technique[%d].name is required", i, j)
			}
		}
	}
	return nil
}

func resolveSwordMaster(def *BookDefinition, masters []models.SwordMaster) (int, *models.SwordMaster) {
	if def.NewSwordMaster != nil {
		m := &models.SwordMaster{
			Name:      def.NewSwordMaster.Name,
			Bio:       def.NewSwordMaster.Bio,
			BirthYear: def.NewSwordMaster.BirthYear,
			DeathYear: def.NewSwordMaster.DeathYear,
		}
		return 0, m
	}
	for _, m := range masters {
		if m.Name == def.SwordMaster {
			return m.ID, nil
		}
	}
	return 0, nil
}

func findBookByTitle(books []models.FightingBook, title string) *models.FightingBook {
	for _, b := range books {
		if b.Title == title {
			return &b
		}
	}
	return nil
}

func maxID[T any](items []T, getID func(T) int) int {
	max := 0
	for _, item := range items {
		if id := getID(item); id > max {
			max = id
		}
	}
	return max
}

func loadJSON[T any](path string) ([]T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return items, nil
}

func writeJSON(path string, data any) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0644)
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
