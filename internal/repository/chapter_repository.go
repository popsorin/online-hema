package repository

import (
	"database/sql"
	"fmt"

	"hema-lessons/internal/models"
)

type ChapterRepository struct {
	db *sql.DB
}

func NewChapterRepository(db *sql.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) ListByFightingBookID(fightingBookID int) ([]models.Chapter, error) {
	query := `
		SELECT id, fighting_book_id, chapter_number, title, description, created_at, updated_at
		FROM chapters
		WHERE fighting_book_id = $1
		ORDER BY chapter_number ASC
	`

	rows, err := r.db.Query(query, fightingBookID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chapters: %w", err)
	}
	defer rows.Close()

	var chapters []models.Chapter
	for rows.Next() {
		var chapter models.Chapter
		err := rows.Scan(
			&chapter.ID,
			&chapter.FightingBookID,
			&chapter.ChapterNumber,
			&chapter.Title,
			&chapter.Description,
			&chapter.CreatedAt,
			&chapter.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chapter: %w", err)
		}
		chapters = append(chapters, chapter)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chapters: %w", err)
	}

	return chapters, nil
}

func (r *ChapterRepository) FightingBookExists(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM fighting_books WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check fighting book existence: %w", err)
	}
	return exists, nil
}
