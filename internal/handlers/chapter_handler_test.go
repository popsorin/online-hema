package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hema-lessons/internal/models"
	"hema-lessons/internal/repository"
	"hema-lessons/internal/testutil"
)

func TestChapterHandler_ListByFightingBook(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	masterIDs := testutil.SeedSwordMasters(t, db)
	bookIDs := testutil.SeedFightingBooks(t, db, masterIDs)
	testutil.SeedChapters(t, db, bookIDs)

	repo := repository.NewChapterRepository(db)
	handler := NewChapterHandler(repo)

	tests := []struct {
		name               string
		path               string
		expectedStatusCode int
		expectedCount      int
	}{
		{
			name:               "get chapters for book with 3 chapters",
			path:               "/api/fighting-books/1/chapters",
			expectedStatusCode: http.StatusOK,
			expectedCount:      3,
		},
		{
			name:               "get chapters for book with 2 chapters",
			path:               "/api/fighting-books/2/chapters",
			expectedStatusCode: http.StatusOK,
			expectedCount:      2,
		},
		{
			name:               "get chapters for book with no chapters",
			path:               "/api/fighting-books/3/chapters",
			expectedStatusCode: http.StatusOK,
			expectedCount:      0,
		},
		{
			name:               "invalid book ID - string",
			path:               "/api/fighting-books/abc/chapters",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid book ID - zero",
			path:               "/api/fighting-books/0/chapters",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid book ID - negative",
			path:               "/api/fighting-books/-1/chapters",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "non-existent book",
			path:               "/api/fighting-books/999/chapters",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ListByFightingBook(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedStatusCode == http.StatusOK {
				var chapters []models.Chapter
				if err := json.NewDecoder(w.Body).Decode(&chapters); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(chapters) != tt.expectedCount {
					t.Errorf("expected %d chapters, got %d", tt.expectedCount, len(chapters))
				}
			}
		})
	}
}

func TestChapterHandler_ListByFightingBook_VerifyOrdering(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	masterIDs := testutil.SeedSwordMasters(t, db)
	bookIDs := testutil.SeedFightingBooks(t, db, masterIDs)
	testutil.SeedChapters(t, db, bookIDs)

	repo := repository.NewChapterRepository(db)
	handler := NewChapterHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/fighting-books/1/chapters", nil)
	w := httptest.NewRecorder()

	handler.ListByFightingBook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var chapters []models.Chapter
	if err := json.NewDecoder(w.Body).Decode(&chapters); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for i := 1; i < len(chapters); i++ {
		if chapters[i].ChapterNumber < chapters[i-1].ChapterNumber {
			t.Errorf("chapters not ordered by chapter_number: chapter %d comes before chapter %d",
				chapters[i-1].ChapterNumber, chapters[i].ChapterNumber)
		}
	}
}

func TestChapterHandler_ListByFightingBook_VerifyFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	masterIDs := testutil.SeedSwordMasters(t, db)
	bookIDs := testutil.SeedFightingBooks(t, db, masterIDs)
	testutil.SeedChapters(t, db, bookIDs)

	repo := repository.NewChapterRepository(db)
	handler := NewChapterHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/fighting-books/1/chapters", nil)
	w := httptest.NewRecorder()

	handler.ListByFightingBook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var chapters []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&chapters); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(chapters) == 0 {
		t.Fatal("expected at least one chapter")
	}

	requiredFields := []string{"id", "fighting_book_id", "chapter_number", "title", "description", "created_at", "updated_at"}
	for _, field := range requiredFields {
		if _, exists := chapters[0][field]; !exists {
			t.Errorf("expected field %s to exist in response", field)
		}
	}
}

func TestChapterHandler_ListByFightingBook_EmptyResult(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	masterIDs := testutil.SeedSwordMasters(t, db)
	testutil.SeedFightingBooks(t, db, masterIDs)

	repo := repository.NewChapterRepository(db)
	handler := NewChapterHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/fighting-books/1/chapters", nil)
	w := httptest.NewRecorder()

	handler.ListByFightingBook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var chapters []models.Chapter
	if err := json.NewDecoder(w.Body).Decode(&chapters); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if chapters != nil && len(chapters) != 0 {
		t.Errorf("expected empty array, got %d chapters", len(chapters))
	}
}
