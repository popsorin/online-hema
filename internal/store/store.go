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

// ResourceWithAuthor combines a Resource with its Author's name.
type ResourceWithAuthor struct {
	models.Resource
	AuthorName string `json:"author_name,omitempty"`
}

// Store holds all application data in memory, loaded from embedded JSON files.
type Store struct {
	authors   map[int]models.Author
	resources map[int]models.Resource
	sections  map[int]models.Section
	items     map[int]models.Item
}

// New creates a Store by parsing the embedded JSON data files.
func New() (*Store, error) {
	s := &Store{}

	if err := s.loadAuthors(); err != nil {
		return nil, fmt.Errorf("loading authors: %w", err)
	}
	if err := s.loadResources(); err != nil {
		return nil, fmt.Errorf("loading resources: %w", err)
	}
	if err := s.loadSections(); err != nil {
		return nil, fmt.Errorf("loading sections: %w", err)
	}
	if err := s.loadItems(); err != nil {
		return nil, fmt.Errorf("loading items: %w", err)
	}

	return s, nil
}

// NewFromData creates a Store from pre-built data (useful for testing).
func NewFromData(
	authors []models.Author,
	resources []models.Resource,
	sections []models.Section,
	items []models.Item,
) *Store {
	s := &Store{
		authors:   make(map[int]models.Author, len(authors)),
		resources: make(map[int]models.Resource, len(resources)),
		sections:  make(map[int]models.Section, len(sections)),
		items:     make(map[int]models.Item, len(items)),
	}

	for _, a := range authors {
		s.authors[a.ID] = a
	}
	for _, r := range resources {
		s.resources[r.ID] = r
	}
	for _, sec := range sections {
		s.sections[sec.ID] = sec
	}
	for _, i := range items {
		s.items[i.ID] = i
	}

	return s
}

// --- Resources ---

// ListResources returns a paginated list of resources with their author names, ordered by title.
func (s *Store) ListResources(params pagination.Params) ([]ResourceWithAuthor, int) {
	totalCount := len(s.resources)
	if totalCount == 0 {
		return nil, 0
	}

	resources := make([]ResourceWithAuthor, 0, totalCount)
	for _, r := range s.resources {
		rwa := ResourceWithAuthor{Resource: r}
		if r.AuthorID != nil {
			if author, ok := s.authors[*r.AuthorID]; ok {
				rwa.AuthorName = author.Name
			}
		}
		resources = append(resources, rwa)
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Title < resources[j].Title
	})

	start := params.Offset
	if start > totalCount {
		start = totalCount
	}
	end := start + params.PageSize
	if end > totalCount {
		end = totalCount
	}

	page := resources[start:end]
	if len(page) == 0 {
		return nil, totalCount
	}

	return page, totalCount
}

// GetResourceByID returns a single resource with its author name, or nil if not found.
func (s *Store) GetResourceByID(id int) *ResourceWithAuthor {
	r, ok := s.resources[id]
	if !ok {
		return nil
	}

	rwa := &ResourceWithAuthor{Resource: r}
	if r.AuthorID != nil {
		if author, ok := s.authors[*r.AuthorID]; ok {
			rwa.AuthorName = author.Name
		}
	}
	return rwa
}

// ResourceExists returns true if a resource with the given ID exists.
func (s *Store) ResourceExists(id int) bool {
	_, ok := s.resources[id]
	return ok
}

// --- Sections ---

// ListRootSectionsByResourceID returns top-level sections (parent_id is nil) for a given resource, ordered by position.
func (s *Store) ListRootSectionsByResourceID(resourceID int) []models.Section {
	var sections []models.Section
	for _, sec := range s.sections {
		if sec.ResourceID == resourceID && sec.ParentID == nil {
			sections = append(sections, sec)
		}
	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Position < sections[j].Position
	})

	return sections
}

// GetSectionByID returns a single section, or nil if not found.
func (s *Store) GetSectionByID(id int) *models.Section {
	sec, ok := s.sections[id]
	if !ok {
		return nil
	}
	return &sec
}

// ListChildSections returns direct child sections of a given parent section, ordered by position.
func (s *Store) ListChildSections(parentID int) []models.Section {
	var sections []models.Section
	for _, sec := range s.sections {
		if sec.ParentID != nil && *sec.ParentID == parentID {
			sections = append(sections, sec)
		}
	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Position < sections[j].Position
	})

	return sections
}

// --- Items ---

// ListItemsBySectionID returns items for a given section, ordered by position.
func (s *Store) ListItemsBySectionID(sectionID int) []models.Item {
	var items []models.Item
	for _, item := range s.items {
		if item.SectionID == sectionID {
			items = append(items, item)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Position < items[j].Position
	})

	return items
}

// --- Data loading ---

func (s *Store) loadAuthors() error {
	var items []models.Author
	if err := loadJSON("data/authors.json", &items); err != nil {
		return err
	}
	s.authors = make(map[int]models.Author, len(items))
	for _, item := range items {
		s.authors[item.ID] = item
	}
	return nil
}

func (s *Store) loadResources() error {
	var items []models.Resource
	if err := loadJSON("data/resources.json", &items); err != nil {
		return err
	}
	s.resources = make(map[int]models.Resource, len(items))
	for _, item := range items {
		s.resources[item.ID] = item
	}
	return nil
}

func (s *Store) loadSections() error {
	var items []models.Section
	if err := loadJSON("data/sections.json", &items); err != nil {
		return err
	}
	s.sections = make(map[int]models.Section, len(items))
	for _, item := range items {
		s.sections[item.ID] = item
	}
	return nil
}

func (s *Store) loadItems() error {
	var items []models.Item
	if err := loadJSON("data/items.json", &items); err != nil {
		return err
	}
	s.items = make(map[int]models.Item, len(items))
	for _, item := range items {
		s.items[item.ID] = item
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
