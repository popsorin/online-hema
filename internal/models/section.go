package models

type Section struct {
	ID         int    `json:"id"`
	ResourceID int    `json:"resource_id"`
	ParentID    *int   `json:"parent_id,omitempty"`
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}