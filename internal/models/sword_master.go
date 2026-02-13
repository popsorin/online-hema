package models

import (
	"time"
)

type SwordMaster struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Bio         string    `json:"bio" db:"bio"`
	BirthYear   *int      `json:"birth_year,omitempty" db:"birth_year"`
	DeathYear   *int      `json:"death_year,omitempty" db:"death_year"`
	ImageURL    *string   `json:"image_url,omitempty" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
