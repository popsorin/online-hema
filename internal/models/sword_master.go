package models

type SwordMaster struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Bio       string  `json:"bio"`
	BirthYear *int    `json:"birth_year,omitempty"`
	DeathYear *int    `json:"death_year,omitempty"`
	ImageURL  *string `json:"image_url,omitempty"`
}
