package store

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"

	"hema-lessons/internal/models"
	"hema-lessons/internal/pagination"
)

//go:embed data
var dataFS embed.FS

// FightingBookWithMaster combines a FightingBook with its SwordMaster's name.
type FightingBookWithMaster struct {
	models.FightingBook
	SwordMasterName string `json:"sword_master_name"`
}

// Store holds all application data in memory, loaded from embedded JSON files.
type Store struct {
	swordMasters map[int]models.SwordMaster
	fightingBooks map[int]models.FightingBook
	chapters      map[int]models.Chapter
	techniques    map[int]models.Technique
}

// New creates a Store by parsing the embedded JSON data files.
func New() (*Store, error) {
	s := &Store{}

	if err := s.loadSwordMasters(); err != nil {
		return nil, fmt.Errorf("loading sword masters: %w", err)
	}
	if err := s.loadFightingBooks(); err != nil {
		return nil, fmt.Errorf("loading fighting books: %w", err)
	}
	if err := s.loadChapters(); err != nil {
		return nil, fmt.Errorf("loading chapters: %w", err)
	}
	if err := s.loadTechniques(); err != nil {
		return nil, fmt.Errorf("loading techniques: %w", err)
	}

	return s, nil
}

// NewFromData creates a Store from pre-built data (useful for testing).
func NewFromData(
	swordMasters []models.SwordMaster,
	fightingBooks []models.FightingBook,
	chapters []models.Chapter,
	techniques []models.Technique,
) *Store {
	s := &Store{
		swordMasters:  make(map[int]models.SwordMaster, len(swordMasters)),
		fightingBooks: make(map[int]models.FightingBook, len(fightingBooks)),
		chapters:      make(map[int]models.Chapter, len(chapters)),
		techniques:    make(map[int]models.Technique, len(techniques)),
	}

	for _, m := range swordMasters {
		s.swordMasters[m.ID] = m
	}
	for _, b := range fightingBooks {
		s.fightingBooks[b.ID] = b
	}
	for _, c := range chapters {
		s.chapters[c.ID] = c
	}
	for _, t := range techniques {
		s.techniques[t.ID] = t
	}

	return s
}

// --- Fighting Books ---

// ListFightingBooks returns a paginated list of fighting books with their sword master names, ordered by title.
func (s *Store) ListFightingBooks(params pagination.Params) ([]FightingBookWithMaster, int) {
	totalCount := len(s.fightingBooks)
	if totalCount == 0 {
		return nil, 0
	}

	// Collect all books into a slice
	books := make([]FightingBookWithMaster, 0, totalCount)
	for _, fb := range s.fightingBooks {
		master := s.swordMasters[fb.SwordMasterID]
		books = append(books, FightingBookWithMaster{
			FightingBook:    fb,
			SwordMasterName: master.Name,
		})
	}

	// Sort by title ASC
	sort.Slice(books, func(i, j int) bool {
		return books[i].Title < books[j].Title
	})

	// Apply pagination
	start := params.Offset
	if start > totalCount {
		start = totalCount
	}
	end := start + params.PageSize
	if end > totalCount {
		end = totalCount
	}

	page := books[start:end]
	if len(page) == 0 {
		return nil, totalCount
	}

	return page, totalCount
}

// GetFightingBookByID returns a single fighting book with its sword master name, or nil if not found.
func (s *Store) GetFightingBookByID(id int) *FightingBookWithMaster {
	fb, ok := s.fightingBooks[id]
	if !ok {
		return nil
	}

	master := s.swordMasters[fb.SwordMasterID]
	return &FightingBookWithMaster{
		FightingBook:    fb,
		SwordMasterName: master.Name,
	}
}

// --- Chapters ---

// FightingBookExists returns true if a fighting book with the given ID exists.
func (s *Store) FightingBookExists(id int) bool {
	_, ok := s.fightingBooks[id]
	return ok
}

// ListChaptersByBookID returns chapters for a given fighting book, ordered by chapter number.
func (s *Store) ListChaptersByBookID(fightingBookID int) []models.Chapter {
	var chapters []models.Chapter
	for _, c := range s.chapters {
		if c.FightingBookID == fightingBookID {
			chapters = append(chapters, c)
		}
	}

	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].ChapterNumber < chapters[j].ChapterNumber
	})

	return chapters
}

// --- Techniques ---

// ChapterExists returns true if a chapter with the given ID exists.
func (s *Store) ChapterExists(id int) bool {
	_, ok := s.chapters[id]
	return ok
}

// ListTechniquesByChapterID returns techniques for a given chapter, ordered by order_in_chapter.
func (s *Store) ListTechniquesByChapterID(chapterID int) []models.Technique {
	var techniques []models.Technique
	for _, t := range s.techniques {
		if t.ChapterID == chapterID {
			techniques = append(techniques, t)
		}
	}

	sort.Slice(techniques, func(i, j int) bool {
		return techniques[i].OrderInChapter < techniques[j].OrderInChapter
	})

	return techniques
}

// --- Data loading ---

func (s *Store) loadSwordMasters() error {
	var items []models.SwordMaster
	if err := loadJSON("data/sword_masters.json", &items); err != nil {
		return err
	}
	s.swordMasters = make(map[int]models.SwordMaster, len(items))
	for _, item := range items {
		s.swordMasters[item.ID] = item
	}
	return nil
}

func (s *Store) loadFightingBooks() error {
	var items []models.FightingBook
	if err := loadJSON("data/fighting_books.json", &items); err != nil {
		return err
	}
	s.fightingBooks = make(map[int]models.FightingBook, len(items))
	for _, item := range items {
		s.fightingBooks[item.ID] = item
	}
	return nil
}

func (s *Store) loadChapters() error {
	var items []models.Chapter
	if err := loadJSON("data/chapters.json", &items); err != nil {
		return err
	}
	s.chapters = make(map[int]models.Chapter, len(items))
	for _, item := range items {
		s.chapters[item.ID] = item
	}
	return nil
}

func (s *Store) loadTechniques() error {
	var items []models.Technique
	if err := loadJSON("data/techniques.json", &items); err != nil {
		return err
	}
	s.techniques = make(map[int]models.Technique, len(items))
	for _, item := range items {
		s.techniques[item.ID] = item
	}
	return nil
}

func loadJSON(path string, dest interface{}) error {
	data, err := dataFS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}
	return nil
}
