package models

import "encoding/json"

type Item struct {
	ID          int             `json:"id"`
	SectionID   int             `json:"section_id"`
	Kind        string          `json:"kind"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Position    int             `json:"position"`
	Attributes  json.RawMessage `json:"attributes,omitempty"`
}
