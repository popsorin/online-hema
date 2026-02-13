package testutil

import (
	"hema-lessons/internal/models"
	"hema-lessons/internal/store"
)

// TestSwordMasters returns a set of sword masters for testing.
func TestSwordMasters() []models.SwordMaster {
	return []models.SwordMaster{
		{ID: 1, Name: "Test Master 1", Bio: "First test master", BirthYear: intPtr(1300), DeathYear: intPtr(1380)},
		{ID: 2, Name: "Test Master 2", Bio: "Second test master", BirthYear: intPtr(1350), DeathYear: intPtr(1420)},
		{ID: 3, Name: "Test Master 3", Bio: "Third test master", BirthYear: intPtr(1400), DeathYear: intPtr(1470)},
	}
}

// TestFightingBooks returns 5 fighting books for testing. Books A and B belong to
// master 1, C and D to master 2, E to master 3.
func TestFightingBooks() []models.FightingBook {
	return []models.FightingBook{
		{ID: 1, SwordMasterID: 1, Title: "Book A", Description: "First test book", PublicationYear: intPtr(1400)},
		{ID: 2, SwordMasterID: 1, Title: "Book B", Description: "Second test book", PublicationYear: intPtr(1410)},
		{ID: 3, SwordMasterID: 2, Title: "Book C", Description: "Third test book", PublicationYear: intPtr(1420)},
		{ID: 4, SwordMasterID: 2, Title: "Book D", Description: "Fourth test book", PublicationYear: intPtr(1430)},
		{ID: 5, SwordMasterID: 3, Title: "Book E", Description: "Fifth test book", PublicationYear: intPtr(1440)},
	}
}

// TestChapters returns chapters for testing. Book 1 has 3 chapters, book 2 has 2 chapters.
func TestChapters() []models.Chapter {
	return []models.Chapter{
		{ID: 1, FightingBookID: 1, ChapterNumber: 1, Title: "Chapter 1", Description: "First chapter of Book A"},
		{ID: 2, FightingBookID: 1, ChapterNumber: 2, Title: "Chapter 2", Description: "Second chapter of Book A"},
		{ID: 3, FightingBookID: 1, ChapterNumber: 3, Title: "Chapter 3", Description: "Third chapter of Book A"},
		{ID: 4, FightingBookID: 2, ChapterNumber: 1, Title: "Introduction", Description: "First chapter of Book B"},
		{ID: 5, FightingBookID: 2, ChapterNumber: 2, Title: "Advanced Techniques", Description: "Second chapter of Book B"},
	}
}

// TestTechniques returns techniques for testing. Chapter 1 has 3 techniques, chapter 2 has 2 techniques.
func TestTechniques() []models.Technique {
	return []models.Technique{
		{ID: 1, ChapterID: 1, Name: "Technique 1", Description: "First technique", Instructions: "Step 1, Step 2", OrderInChapter: 1},
		{ID: 2, ChapterID: 1, Name: "Technique 2", Description: "Second technique", Instructions: "Step A, Step B", OrderInChapter: 2},
		{ID: 3, ChapterID: 1, Name: "Technique 3", Description: "Third technique", Instructions: "Instructions here", OrderInChapter: 3},
		{ID: 4, ChapterID: 2, Name: "Basic Move", Description: "A basic move", Instructions: "Do this", OrderInChapter: 1},
		{ID: 5, ChapterID: 2, Name: "Advanced Move", Description: "An advanced move", Instructions: "Then do this", OrderInChapter: 2},
	}
}

// NewTestStore creates a Store populated with all test data.
func NewTestStore() *store.Store {
	return store.NewFromData(
		TestSwordMasters(),
		TestFightingBooks(),
		TestChapters(),
		TestTechniques(),
	)
}

// NewEmptyStore creates an empty Store for testing empty-result scenarios.
func NewEmptyStore() *store.Store {
	return store.NewFromData(nil, nil, nil, nil)
}

// NewStoreWithMastersAndBooks creates a Store with only masters and books (no chapters/techniques).
func NewStoreWithMastersAndBooks() *store.Store {
	return store.NewFromData(
		TestSwordMasters(),
		TestFightingBooks(),
		nil,
		nil,
	)
}

// NewStoreWithMastersBooksAndChapters creates a Store with masters, books, and chapters (no techniques).
func NewStoreWithMastersBooksAndChapters() *store.Store {
	return store.NewFromData(
		TestSwordMasters(),
		TestFightingBooks(),
		TestChapters(),
		nil,
	)
}

func intPtr(v int) *int {
	return &v
}
