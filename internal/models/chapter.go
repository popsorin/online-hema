package models

type Chapter struct {
	ID             int    `json:"id"`
	FightingBookID int    `json:"fighting_book_id"`
	ChapterNumber  int    `json:"chapter_number"`
	Title          string `json:"title"`
	Description    string `json:"description"`
}
