package repository

import (
	"database/sql"
	"fmt"

	"hema-lessons/internal/models"
)

type TechniqueRepository struct {
	db *sql.DB
}

func NewTechniqueRepository(db *sql.DB) *TechniqueRepository {
	return &TechniqueRepository{db: db}
}

func (r *TechniqueRepository) ListByChapterID(chapterID int) ([]models.Technique, error) {
	query := `
		SELECT id, chapter_id, name, description, instructions, video_url, thumbnail_url, order_in_chapter, created_at, updated_at
		FROM techniques
		WHERE chapter_id = $1
		ORDER BY order_in_chapter ASC
	`

	rows, err := r.db.Query(query, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to query techniques: %w", err)
	}
	defer rows.Close()

	var techniques []models.Technique
	for rows.Next() {
		var technique models.Technique
		err := rows.Scan(
			&technique.ID,
			&technique.ChapterID,
			&technique.Name,
			&technique.Description,
			&technique.Instructions,
			&technique.VideoURL,
			&technique.ThumbnailURL,
			&technique.OrderInChapter,
			&technique.CreatedAt,
			&technique.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan technique: %w", err)
		}
		techniques = append(techniques, technique)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating techniques: %w", err)
	}

	return techniques, nil
}

func (r *TechniqueRepository) ChapterExists(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM chapters WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check chapter existence: %w", err)
	}
	return exists, nil
}
