package testutil

import (
	"encoding/json"

	"hema-lessons/internal/models"
	"hema-lessons/internal/store"
)

// TestAuthors returns a set of authors for testing.
func TestAuthors() []models.Author {
	return []models.Author{
		{ID: 1, Name: "Test Author 1", Bio: "First test author", BirthYear: intPtr(1300), DeathYear: intPtr(1380)},
		{ID: 2, Name: "Test Author 2", Bio: "Second test author", BirthYear: intPtr(1350), DeathYear: intPtr(1420)},
		{ID: 3, Name: "Test Author 3", Bio: "Third test author", BirthYear: intPtr(1400), DeathYear: intPtr(1470)},
	}
}

// TestResources returns 5 resources for testing. Resources A and B belong to author 1,
// C and D to author 2, E to author 3.
func TestResources() []models.Resource {
	return []models.Resource{
		{ID: 1, AuthorID: intPtr(1), Title: "Book A", Description: "First test book", PublicationYear: intPtr(1400)},
		{ID: 2, AuthorID: intPtr(1), Title: "Book B", Description: "Second test book", PublicationYear: intPtr(1410)},
		{ID: 3, AuthorID: intPtr(2), Title: "Book C", Description: "Third test book", PublicationYear: intPtr(1420)},
		{ID: 4, AuthorID: intPtr(2), Title: "Book D", Description: "Fourth test book", PublicationYear: intPtr(1430)},
		{ID: 5, AuthorID: intPtr(3), Title: "Book E", Description: "Fifth test book", PublicationYear: intPtr(1440)},
	}
}

// TestSections returns sections for testing.
// Resource 1 has 3 root sections; resource 2 has 2 root sections.
// Section 1 also has 1 child section (id=6) to test arbitrary nesting.
func TestSections() []models.Section {
	return []models.Section{
		{ID: 1, ResourceID: 1, ParentID: nil, Kind: "chapter", Title: "Chapter 1", Description: "First chapter of Book A", Position: 1},
		{ID: 2, ResourceID: 1, ParentID: nil, Kind: "chapter", Title: "Chapter 2", Description: "Second chapter of Book A", Position: 2},
		{ID: 3, ResourceID: 1, ParentID: nil, Kind: "chapter", Title: "Chapter 3", Description: "Third chapter of Book A", Position: 3},
		{ID: 4, ResourceID: 2, ParentID: nil, Kind: "chapter", Title: "Introduction", Description: "First chapter of Book B", Position: 1},
		{ID: 5, ResourceID: 2, ParentID: nil, Kind: "chapter", Title: "Advanced Techniques", Description: "Second chapter of Book B", Position: 2},
		{ID: 6, ResourceID: 1, ParentID: intPtr(1), Kind: "sub-chapter", Title: "Sub-section of Chapter 1", Description: "A nested section", Position: 1},
	}
}

// TestItems returns items for testing.
// Section 1 has 3 items; section 2 has 2 items.
func TestItems() []models.Item {
	attrs := json.RawMessage(`{"instructions":"Step 1, Step 2"}`)
	return []models.Item{
		{ID: 1, SectionID: 1, Kind: "technique", Title: "Technique 1", Description: "First technique", Position: 1, Attributes: attrs},
		{ID: 2, SectionID: 1, Kind: "technique", Title: "Technique 2", Description: "Second technique", Position: 2, Attributes: attrs},
		{ID: 3, SectionID: 1, Kind: "technique", Title: "Technique 3", Description: "Third technique", Position: 3, Attributes: attrs},
		{ID: 4, SectionID: 2, Kind: "technique", Title: "Basic Move", Description: "A basic move", Position: 1, Attributes: attrs},
		{ID: 5, SectionID: 2, Kind: "technique", Title: "Advanced Move", Description: "An advanced move", Position: 2, Attributes: attrs},
	}
}

// NewTestStore creates a Store populated with all test data.
func NewTestStore() *store.Store {
	return store.NewFromData(
		TestAuthors(),
		TestResources(),
		TestSections(),
		TestItems(),
	)
}

// NewEmptyStore creates an empty Store for testing empty-result scenarios.
func NewEmptyStore() *store.Store {
	return store.NewFromData(nil, nil, nil, nil)
}

// NewStoreWithAuthorsAndResources creates a Store with only authors and resources (no sections/items).
func NewStoreWithAuthorsAndResources() *store.Store {
	return store.NewFromData(
		TestAuthors(),
		TestResources(),
		nil,
		nil,
	)
}

// NewStoreWithAuthorsResourcesSections creates a Store with authors, resources, and sections (no items).
func NewStoreWithAuthorsResourcesSections() *store.Store {
	return store.NewFromData(
		TestAuthors(),
		TestResources(),
		TestSections(),
		nil,
	)
}

func intPtr(v int) *int {
	return &v
}
