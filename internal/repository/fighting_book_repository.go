package repository

import (
	"database/sql"
	"fmt"

	"hema-lessons/internal/models"
	"hema-lessons/internal/pagination"
)

type FightingBookRepository struct {
	db *sql.DB
}

func NewFightingBookRepository(db *sql.DB) *FightingBookRepository {
	return &FightingBookRepository{db: db}
}

type FightingBookWithMaster struct {
	models.FightingBook
	SwordMasterName string `json:"sword_master_name"`
}

func (r *FightingBookRepository) List(params pagination.Params) ([]FightingBookWithMaster, int, error) {
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM fighting_books"
	if err := r.db.QueryRow(countQuery).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count fighting books: %w", err)
	}

	query := `
		SELECT
			fb.id, fb.sword_master_id, fb.title, fb.description,
			fb.publication_year, fb.cover_image_url, fb.created_at, fb.updated_at,
			sm.name as sword_master_name
		FROM fighting_books fb
		INNER JOIN sword_masters sm ON fb.sword_master_id = sm.id
		ORDER BY fb.title ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, params.PageSize, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query fighting books: %w", err)
	}
	defer rows.Close()

	var books []FightingBookWithMaster
	for rows.Next() {
		var book FightingBookWithMaster
		err := rows.Scan(
			&book.ID,
			&book.SwordMasterID,
			&book.Title,
			&book.Description,
			&book.PublicationYear,
			&book.CoverImageURL,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.SwordMasterName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fighting book: %w", err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating fighting books: %w", err)
	}

	return books, totalCount, nil
}

func (r *FightingBookRepository) GetByID(id int) (*FightingBookWithMaster, error) {
	query := `
		SELECT 
			fb.id, fb.sword_master_id, fb.title, fb.description, 
			fb.publication_year, fb.cover_image_url, fb.created_at, fb.updated_at,
			sm.name as sword_master_name
		FROM fighting_books fb
		INNER JOIN sword_masters sm ON fb.sword_master_id = sm.id
		WHERE fb.id = $1
	`

	var book FightingBookWithMaster
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.SwordMasterID,
		&book.Title,
		&book.Description,
		&book.PublicationYear,
		&book.CoverImageURL,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.SwordMasterName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fighting book: %w", err)
	}

	return &book, nil
}
