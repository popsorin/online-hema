package models

type Resource struct {
	ID              int     `json:"id"`
	AuthorID        *int    `json:"author_id,omitempty"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	PublicationYear *int    `json:"publication_year,omitempty"`
	CoverImageURL   *string `json:"cover_image_url,omitempty"`
}
