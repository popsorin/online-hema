package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"hema-lessons/internal/pagination"
	"hema-lessons/internal/store"
	"hema-lessons/internal/testutil"
)

func TestResourceHandler_List(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewResourceHandler(s)

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
			req := httptest.NewRequest(http.MethodGet, "/api/resources"+tt.queryParams, nil)
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

			if tt.expectedDataCount == 0 && (response.Data == nil || len(data) == 0) {
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
				requiredFields := []string{"id", "author_id", "title", "description", "author_name"}
				for _, field := range requiredFields {
					if _, exists := firstItem[field]; !exists {
						t.Errorf("expected field %s to exist in response", field)
					}
				}
			}
		})
	}
}

func TestResourceHandler_List_EmptyDatabase(t *testing.T) {
	s := testutil.NewEmptyStore()
	handler := NewResourceHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/resources", nil)
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

func TestResourceHandler_List_VerifyOrdering(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewResourceHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/resources", nil)
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
		resource := item.(map[string]interface{})
		title := resource["title"].(string)
		if i > 0 && title < prevTitle {
			t.Errorf("resources not ordered by title: %q came after %q", title, prevTitle)
		}
		prevTitle = title
	}
}

func TestResourceHandler_Get(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewResourceHandler(s)

	tests := []struct {
		name               string
		resourceID         int
		expectedStatusCode int
		expectResource     bool
		expectedTitle      string
	}{
		{
			name:               "existing resource - first resource",
			resourceID:         1,
			expectedStatusCode: http.StatusOK,
			expectResource:     true,
			expectedTitle:      "Book A",
		},
		{
			name:               "existing resource - last resource",
			resourceID:         5,
			expectedStatusCode: http.StatusOK,
			expectResource:     true,
			expectedTitle:      "Book E",
		},
		{
			name:               "non-existing resource",
			resourceID:         99999,
			expectedStatusCode: http.StatusNotFound,
			expectResource:     false,
		},
		{
			name:               "invalid resource ID",
			resourceID:         -1,
			expectedStatusCode: http.StatusBadRequest,
			expectResource:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/resources/%d", tt.resourceID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handler.Get(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if !tt.expectResource {
				return
			}

			var resource store.ResourceWithAuthor
			if err := json.NewDecoder(w.Body).Decode(&resource); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resource.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, resource.Title)
			}
			if resource.ID == 0 {
				t.Error("expected ID to be set")
			}
			if resource.AuthorName == "" {
				t.Error("expected AuthorName to be set")
			}
		})
	}
}
