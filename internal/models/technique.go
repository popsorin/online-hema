package models

type Technique struct {
	ID             int     `json:"id"`
	ChapterID      int     `json:"chapter_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Instructions   string  `json:"instructions"`
	VideoURL       *string `json:"video_url,omitempty"`
	ThumbnailURL   *string `json:"thumbnail_url,omitempty"`
	OrderInChapter int     `json:"order_in_chapter"`
}
