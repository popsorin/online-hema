package models

import (
	"time"
)

type Chapter struct {
	ID             int       `json:"id" db:"id"`
	FightingBookID int       `json:"fighting_book_id" db:"fighting_book_id"`
	ChapterNumber  int       `json:"chapter_number" db:"chapter_number"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
