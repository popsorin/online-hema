package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hema-lessons/internal/models"
	"hema-lessons/internal/testutil"
)

func TestTechniqueHandler_ListByChapter(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewTechniqueHandler(s)

	tests := []struct {
		name               string
		path               string
		expectedStatusCode int
		expectedCount      int
	}{
		{
			name:               "get techniques for chapter with 3 techniques",
			path:               "/api/chapters/1/techniques",
			expectedStatusCode: http.StatusOK,
			expectedCount:      3,
		},
		{
			name:               "get techniques for chapter with 2 techniques",
			path:               "/api/chapters/2/techniques",
			expectedStatusCode: http.StatusOK,
			expectedCount:      2,
		},
		{
			name:               "get techniques for chapter with no techniques",
			path:               "/api/chapters/3/techniques",
			expectedStatusCode: http.StatusOK,
			expectedCount:      0,
		},
		{
			name:               "invalid chapter ID - string",
			path:               "/api/chapters/abc/techniques",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid chapter ID - zero",
			path:               "/api/chapters/0/techniques",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid chapter ID - negative",
			path:               "/api/chapters/-1/techniques",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "non-existent chapter",
			path:               "/api/chapters/999/techniques",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ListByChapter(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedStatusCode == http.StatusOK {
				var techniques []models.Technique
				if err := json.NewDecoder(w.Body).Decode(&techniques); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(techniques) != tt.expectedCount {
					t.Errorf("expected %d techniques, got %d", tt.expectedCount, len(techniques))
				}
			}
		})
	}
}

func TestTechniqueHandler_ListByChapter_VerifyOrdering(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewTechniqueHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/chapters/1/techniques", nil)
	w := httptest.NewRecorder()

	handler.ListByChapter(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var techniques []models.Technique
	if err := json.NewDecoder(w.Body).Decode(&techniques); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for i := 1; i < len(techniques); i++ {
		if techniques[i].OrderInChapter < techniques[i-1].OrderInChapter {
			t.Errorf("techniques not ordered by order_in_chapter: %d comes before %d",
				techniques[i-1].OrderInChapter, techniques[i].OrderInChapter)
		}
	}
}

func TestTechniqueHandler_ListByChapter_VerifyFields(t *testing.T) {
	s := testutil.NewTestStore()
	handler := NewTechniqueHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/chapters/1/techniques", nil)
	w := httptest.NewRecorder()

	handler.ListByChapter(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var techniques []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&techniques); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(techniques) == 0 {
		t.Fatal("expected at least one technique")
	}

	requiredFields := []string{"id", "chapter_id", "name", "description", "instructions", "order_in_chapter"}
	for _, field := range requiredFields {
		if _, exists := techniques[0][field]; !exists {
			t.Errorf("expected field %s to exist in response", field)
		}
	}
}

func TestTechniqueHandler_ListByChapter_EmptyResult(t *testing.T) {
	s := testutil.NewStoreWithMastersBooksAndChapters()
	handler := NewTechniqueHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/chapters/1/techniques", nil)
	w := httptest.NewRecorder()

	handler.ListByChapter(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var techniques []models.Technique
	if err := json.NewDecoder(w.Body).Decode(&techniques); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if techniques != nil && len(techniques) != 0 {
		t.Errorf("expected empty array, got %d techniques", len(techniques))
	}
}
