package models

import (
	"time"
)

type Technique struct {
	ID             int       `json:"id" db:"id"`
	ChapterID      int       `json:"chapter_id" db:"chapter_id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	Instructions   string    `json:"instructions" db:"instructions"`
	VideoURL       *string   `json:"video_url,omitempty" db:"video_url"` // Only for subscribers
	ThumbnailURL   *string   `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	OrderInChapter int       `json:"order_in_chapter" db:"order_in_chapter"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
