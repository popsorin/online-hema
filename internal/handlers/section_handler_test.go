package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hema-lessons/internal/models"
	"hema-lessons/internal/testutil"
)

func TestSectionHandler_ListByBook(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewSectionHandler(s)

	tests := []struct {
		name               string
		path               string
		expectedStatusCode int
		expectedCount      int
	}{
		{
			name:               "resource with 3 root sections",
			path:               "/api/resources/1/sections",
			expectedStatusCode: http.StatusOK,
			expectedCount:      3,
		},
		{
			name:               "resource with 2 root sections",
			path:               "/api/resources/2/sections",
			expectedStatusCode: http.StatusOK,
			expectedCount:      2,
		},
		{
			name:               "resource with no sections",
			path:               "/api/resources/3/sections",
			expectedStatusCode: http.StatusOK,
			expectedCount:      0,
		},
		{
			name:               "invalid resource ID - string",
			path:               "/api/resources/abc/sections",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid resource ID - zero",
			path:               "/api/resources/0/sections",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid resource ID - negative",
			path:               "/api/resources/-1/sections",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "non-existent resource",
			path:               "/api/resources/999/sections",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ListByBook(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedStatusCode == http.StatusOK {
				var sections []models.Section
				if err := json.NewDecoder(w.Body).Decode(&sections); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(sections) != tt.expectedCount {
					t.Errorf("expected %d sections, got %d", tt.expectedCount, len(sections))
				}
			}
		})
	}
}

func TestSectionHandler_ListByBook_VerifyOrdering(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewSectionHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/resources/1/sections", nil)
	w := httptest.NewRecorder()

	handler.ListByBook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var sections []models.Section
	if err := json.NewDecoder(w.Body).Decode(&sections); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for i := 1; i < len(sections); i++ {
		if sections[i].Position < sections[i-1].Position {
			t.Errorf("sections not ordered by position: %d comes before %d",
				sections[i-1].Position, sections[i].Position)
		}
	}
}

func TestSectionHandler_ListByBook_VerifyFields(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewSectionHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/resources/1/sections", nil)
	w := httptest.NewRecorder()

	handler.ListByBook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var sections []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&sections); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(sections) == 0 {
		t.Fatal("expected at least one section")
	}

	requiredFields := []string{"id", "resource_id", "kind", "title", "description", "position"}
	for _, field := range requiredFields {
		if _, exists := sections[0][field]; !exists {
			t.Errorf("expected field %s to exist in response", field)
		}
	}
}

func TestSectionHandler_ListByBook_OnlyReturnsRootSections(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewSectionHandler(s)

	// Resource 1 has sections 1, 2, 3 at root and section 6 as child of section 1.
	// Only 3 root sections should be returned.
	req := httptest.NewRequest(http.MethodGet, "/api/resources/1/sections", nil)
	w := httptest.NewRecorder()

	handler.ListByBook(w, req)

	var sections []models.Section
	if err := json.NewDecoder(w.Body).Decode(&sections); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(sections) != 3 {
		t.Errorf("expected 3 root sections, got %d", len(sections))
	}

	for _, sec := range sections {
		if sec.ParentID != nil {
			t.Errorf("expected all returned sections to be root sections, but section %d has parent_id %d", sec.ID, *sec.ParentID)
		}
	}
}

func TestSectionHandler_ListByBook_EmptyResult(t *testing.T) {
	s := testutil.NewStoreWithAuthorsAndResources()
	handler := NewSectionHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/resources/1/sections", nil)
	w := httptest.NewRecorder()

	handler.ListByBook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var sections []models.Section
	if err := json.NewDecoder(w.Body).Decode(&sections); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if sections != nil && len(sections) != 0 {
		t.Errorf("expected empty result, got %d sections", len(sections))
	}
}

func TestSectionHandler_Get(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewSectionHandler(s)

	tests := []struct {
		name               string
		path               string
		expectedStatusCode int
		expectedTitle      string
	}{
		{
			name:               "existing section",
			path:               "/api/sections/1",
			expectedStatusCode: http.StatusOK,
			expectedTitle:      "Chapter 1",
		},
		{
			name:               "nested section",
			path:               "/api/sections/6",
			expectedStatusCode: http.StatusOK,
			expectedTitle:      "Sub-section of Chapter 1",
		},
		{
			name:               "non-existent section",
			path:               "/api/sections/999",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "invalid section ID",
			path:               "/api/sections/abc",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.Get(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedTitle != "" {
				var sec models.Section
				if err := json.NewDecoder(w.Body).Decode(&sec); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if sec.Title != tt.expectedTitle {
					t.Errorf("expected title %q, got %q", tt.expectedTitle, sec.Title)
				}
			}
		})
	}
}

func TestSectionHandler_ListChildren(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewSectionHandler(s)

	tests := []struct {
		name               string
		path               string
		expectedStatusCode int
		expectedCount      int
	}{
		{
			name:               "section with 1 child",
			path:               "/api/sections/1/sections",
			expectedStatusCode: http.StatusOK,
			expectedCount:      1,
		},
		{
			name:               "section with no children",
			path:               "/api/sections/2/sections",
			expectedStatusCode: http.StatusOK,
			expectedCount:      0,
		},
		{
			name:               "non-existent parent section",
			path:               "/api/sections/999/sections",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "invalid section ID",
			path:               "/api/sections/abc/sections",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ListChildren(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedStatusCode == http.StatusOK {
				var sections []models.Section
				if err := json.NewDecoder(w.Body).Decode(&sections); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(sections) != tt.expectedCount {
					t.Errorf("expected %d sections, got %d", tt.expectedCount, len(sections))
				}
			}
		})
	}
}
