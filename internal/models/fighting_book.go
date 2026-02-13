package models

type FightingBook struct {
	ID              int     `json:"id"`
	SwordMasterID   int     `json:"sword_master_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	PublicationYear *int    `json:"publication_year,omitempty"`
	CoverImageURL   *string `json:"cover_image_url,omitempty"`
}
