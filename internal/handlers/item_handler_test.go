package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hema-lessons/internal/models"
	"hema-lessons/internal/testutil"
)

func TestItemHandler_ListBySection(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewItemHandler(s)

	tests := []struct {
		name               string
		path               string
		expectedStatusCode int
		expectedCount      int
	}{
		{
			name:               "section with 3 items",
			path:               "/api/sections/1/items",
			expectedStatusCode: http.StatusOK,
			expectedCount:      3,
		},
		{
			name:               "section with 2 items",
			path:               "/api/sections/2/items",
			expectedStatusCode: http.StatusOK,
			expectedCount:      2,
		},
		{
			name:               "section with no items",
			path:               "/api/sections/3/items",
			expectedStatusCode: http.StatusOK,
			expectedCount:      0,
		},
		{
			name:               "invalid section ID - string",
			path:               "/api/sections/abc/items",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid section ID - zero",
			path:               "/api/sections/0/items",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid section ID - negative",
			path:               "/api/sections/-1/items",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "non-existent section",
			path:               "/api/sections/999/items",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ListBySection(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedStatusCode == http.StatusOK {
				var items []models.Item
				if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(items) != tt.expectedCount {
					t.Errorf("expected %d items, got %d", tt.expectedCount, len(items))
				}
			}
		})
	}
}

func TestItemHandler_ListBySection_VerifyOrdering(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewItemHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/sections/1/items", nil)
	w := httptest.NewRecorder()

	handler.ListBySection(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var items []models.Item
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for i := 1; i < len(items); i++ {
		if items[i].Position < items[i-1].Position {
			t.Errorf("items not ordered by position: %d comes before %d",
				items[i-1].Position, items[i].Position)
		}
	}
}

func TestItemHandler_ListBySection_VerifyFields(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewItemHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/sections/1/items", nil)
	w := httptest.NewRecorder()

	handler.ListBySection(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var items []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(items) == 0 {
		t.Fatal("expected at least one item")
	}

	requiredFields := []string{"id", "section_id", "kind", "title", "description", "position"}
	for _, field := range requiredFields {
		if _, exists := items[0][field]; !exists {
			t.Errorf("expected field %s to exist in response", field)
		}
	}

	// Verify attributes is present and is an object
	if _, exists := items[0]["attributes"]; !exists {
		t.Error("expected field attributes to exist in response")
	}
}

func TestItemHandler_ListBySection_EmptyResult(t *testing.T) {
	s := testutil.NewStoreWithAuthorsResourcesSections()
	handler := NewItemHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/sections/1/items", nil)
	w := httptest.NewRecorder()

	handler.ListBySection(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var items []models.Item
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if items != nil && len(items) != 0 {
		t.Errorf("expected empty result, got %d items", len(items))
	}
}
