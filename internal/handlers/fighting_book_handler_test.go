package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"hema-lessons/internal/pagination"
	"hema-lessons/internal/repository"
	"hema-lessons/internal/testutil"
)

func TestFightingBookHandler_List(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	masterIDs := testutil.SeedSwordMasters(t, db)
	testutil.SeedFightingBooks(t, db, masterIDs)

	repo := repository.NewFightingBookRepository(db)
	handler := NewFightingBookHandler(repo)

	tests := []struct {
		name               string
		queryParams        string
		expectedStatusCode int
		expectedPage       int
		expectedPageSize   int
		expectedTotalCount int
		expectedTotalPages int
		expectedDataCount  int
	}{
		{
			name:               "default pagination",
			queryParams:        "",
			expectedStatusCode: http.StatusOK,
			expectedPage:       1,
			expectedPageSize:   20,
			expectedTotalCount: 5,
			expectedTotalPages: 1,
			expectedDataCount:  5,
		},
		{
			name:               "custom page size",
			queryParams:        "?page_size=2",
			expectedStatusCode: http.StatusOK,
			expectedPage:       1,
			expectedPageSize:   2,
			expectedTotalCount: 5,
			expectedTotalPages: 3,
			expectedDataCount:  2,
		},
		{
			name:               "page 2 with page size 2",
			queryParams:        "?page=2&page_size=2",
			expectedStatusCode: http.StatusOK,
			expectedPage:       2,
			expectedPageSize:   2,
			expectedTotalCount: 5,
			expectedTotalPages: 3,
			expectedDataCount:  2,
		},
		{
			name:               "last page",
			queryParams:        "?page=3&page_size=2",
			expectedStatusCode: http.StatusOK,
			expectedPage:       3,
			expectedPageSize:   2,
			expectedTotalCount: 5,
			expectedTotalPages: 3,
			expectedDataCount:  1,
		},
		{
			name:               "page beyond available data",
			queryParams:        "?page=10&page_size=2",
			expectedStatusCode: http.StatusOK,
			expectedPage:       10,
			expectedPageSize:   2,
			expectedTotalCount: 5,
			expectedTotalPages: 3,
			expectedDataCount:  0,
		},
		{
			name:               "max page size constraint",
			queryParams:        "?page_size=200",
			expectedStatusCode: http.StatusOK,
			expectedPage:       1,
			expectedPageSize:   100,
			expectedTotalCount: 5,
			expectedTotalPages: 1,
			expectedDataCount:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/fighting-books"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.List(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			var response pagination.Response
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Page != tt.expectedPage {
				t.Errorf("expected page %d, got %d", tt.expectedPage, response.Page)
			}

			if response.PageSize != tt.expectedPageSize {
				t.Errorf("expected page_size %d, got %d", tt.expectedPageSize, response.PageSize)
			}

			if response.TotalCount != tt.expectedTotalCount {
				t.Errorf("expected total_count %d, got %d", tt.expectedTotalCount, response.TotalCount)
			}

			if response.TotalPages != tt.expectedTotalPages {
				t.Errorf("expected total_pages %d, got %d", tt.expectedTotalPages, response.TotalPages)
			}

			data, ok := response.Data.([]interface{})
			if tt.expectedDataCount > 0 && !ok {
				t.Fatalf("expected data to be an array")
			}

			if tt.expectedDataCount == 0 && response.Data == nil {
				// This is fine for empty results
				return
			}

			if len(data) != tt.expectedDataCount {
				t.Errorf("expected %d items in data, got %d", tt.expectedDataCount, len(data))
			}

			if len(data) > 0 {
				firstItem, ok := data[0].(map[string]interface{})
				if !ok {
					t.Fatalf("expected first item to be a map")
				}

				requiredFields := []string{"id", "sword_master_id", "title", "description", "sword_master_name"}
				for _, field := range requiredFields {
					if _, exists := firstItem[field]; !exists {
						t.Errorf("expected field %s to exist in response", field)
					}
				}
			}
		})
	}
}

func TestFightingBookHandler_List_EmptyDatabase(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := repository.NewFightingBookRepository(db)
	handler := NewFightingBookHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/fighting-books", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response pagination.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.TotalCount != 0 {
		t.Errorf("expected total_count 0, got %d", response.TotalCount)
	}

	if response.Data != nil {
		t.Errorf("expected data to be nil for empty result, got %v", response.Data)
	}
}

func TestFightingBookHandler_List_VerifyOrdering(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	masterIDs := testutil.SeedSwordMasters(t, db)
	testutil.SeedFightingBooks(t, db, masterIDs)

	repo := repository.NewFightingBookRepository(db)
	handler := NewFightingBookHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/fighting-books", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	var response pagination.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data := response.Data.([]interface{})
	if len(data) < 2 {
		t.Fatal("not enough data to verify ordering")
	}

	var prevTitle string
	for i, item := range data {
		book := item.(map[string]interface{})
		title := book["title"].(string)

		if i > 0 && title < prevTitle {
			t.Errorf("books are not properly ordered by title: %s came after %s", title, prevTitle)
		}

		prevTitle = title
	}
}

func TestFightingBookHandler_Get(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	masterIDs := testutil.SeedSwordMasters(t, db)
	bookIDs := testutil.SeedFightingBooks(t, db, masterIDs)

	repo := repository.NewFightingBookRepository(db)
	handler := NewFightingBookHandler(repo)

	tests := []struct {
		name               string
		bookID             int
		expectedStatusCode int
		expectBook         bool
		expectedTitle      string
	}{
		{
			name:               "existing book - first book",
			bookID:             bookIDs[0],
			expectedStatusCode: http.StatusOK,
			expectBook:         true,
			expectedTitle:      "Book A",
		},
		{
			name:               "existing book - last book",
			bookID:             bookIDs[len(bookIDs)-1],
			expectedStatusCode: http.StatusOK,
			expectBook:         true,
			expectedTitle:      "Book E",
		},
		{
			name:               "non-existing book",
			bookID:             99999,
			expectedStatusCode: http.StatusNotFound,
			expectBook:         false,
		},
		{
			name:               "invalid book ID",
			bookID:             -1,
			expectedStatusCode: http.StatusBadRequest,
			expectBook:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/fighting-books/%d", tt.bookID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handler.Get(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if !tt.expectBook {
				return
			}

			var book repository.FightingBookWithMaster
			if err := json.NewDecoder(w.Body).Decode(&book); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if book.Title != tt.expectedTitle {
				t.Errorf("expected title %s, got %s", tt.expectedTitle, book.Title)
			}

			if book.ID == 0 {
				t.Error("expected ID to be set")
			}
			if book.SwordMasterName == "" {
				t.Error("expected SwordMasterName to be set")
			}
			if book.Title == "" {
				t.Error("expected Title to be set")
			}
		})
	}
}
