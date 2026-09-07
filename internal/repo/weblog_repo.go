package repo

import (
	"database/sql"
	"fmt"

	"github.com/NooraPanahi/Weblog-Application.git/internal/model"
)

type WeblogRepo struct {
	db *sql.DB
}


func NewWeblogRepo(db *sql.DB) *WeblogRepo {
	return &WeblogRepo{db: db}
}

func (we *WeblogRepo) Create(weblog *model.Weblog) error {
	query := `INSERT INTO weblogs (title, content, image, author_id, privacy)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, created_at`

	err := we.db.QueryRow(query,weblog.Title, weblog.Content, weblog.AuthorID, weblog.Privacy).Scan(&weblog.ID, &weblog.CreatedAt)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *WeblogRepo) FindWeblogByUserID(id int64) (*model.Weblog, error) {
	query := `SELECT id, title, content, image, author_id, privacy, created_at
			  FROM weblogs WHERE id = $1`

	weblog := &model.Weblog{}

	err := r.db.QueryRow(query, id).Scan( &weblog.ID, &weblog.Title, &weblog.Content, &weblog.Image,
		&weblog.AuthorID, &weblog.Privacy, &weblog.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("find weblog by id: %w", err)
	}

	return weblog, nil
}
