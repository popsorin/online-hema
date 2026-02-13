package models

import (
	"time"
)

type FightingBook struct {
	ID              int       `json:"id" db:"id"`
	SwordMasterID   int       `json:"sword_master_id" db:"sword_master_id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	PublicationYear *int      `json:"publication_year,omitempty" db:"publication_year"`
	CoverImageURL   *string   `json:"cover_image_url,omitempty" db:"cover_image_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
