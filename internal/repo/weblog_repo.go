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

	err := we.db.QueryRow(query, weblog.Title, weblog.Content,weblog.Image, weblog.AuthorID, weblog.Privacy).Scan(&weblog.ID, &weblog.CreatedAt)

	if err != nil {
		return fmt.Errorf("create weblog: %w", err)
	}
	return nil
}

func (r *WeblogRepo) FindWeblogByID(id int64) (*model.Weblog, error) {
	query := `SELECT id, title, content, image, author_id, privacy, created_at
			  FROM weblogs WHERE id = $1`

	weblog := &model.Weblog{}

	err := r.db.QueryRow(query, id).Scan(&weblog.ID, &weblog.Title, &weblog.Content, &weblog.Image,
		&weblog.AuthorID, &weblog.Privacy, &weblog.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("find weblog by id: %w", err)
	}

	return weblog, nil
}

func (r *WeblogRepo) ListVisible(userID int64) ([]*model.Weblog, error) {
	query := `SELECT w.id, w.title, w.content, w.image, w.author_id, w.privacy, w.created_at
			  FROM weblogs w 
			  LEFT JOIN weblog_shares ws
			  ON ws.weblog_id = w.id AND ws.user_id = $1
			  WHERE w.privacy = 'public' OR w.author_id = $1
			  OR ws.user_id IS NOT NULL
			  ORDER BY w.created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return  nil, err
	}

	defer rows.Close()

	var weblogs []*model.Weblog

	for rows.Next() {
		weblog := &model.Weblog{}
		err := rows.Scan(&weblog.ID, &weblog.Title, &weblog.Content, &weblog.Image, &weblog.AuthorID, &weblog.Privacy, &weblog.CreatedAt)

		if err != nil {
			return nil, err
		}
		weblogs = append(weblogs, weblog)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return weblogs, nil
}

func (r *WeblogRepo) FindVisibleByID (id, userID int64) (*model.Weblog, error) {
	query := `SELECT w.id, w.title, w.content, w.image, w.author_id, w.privacy, w.created_at FROM weblogs w
			  LEFT JOIN weblog_shares ws 
			  ON ws.weblog_id = w.id 
			  AND ws.user_id = $2
			  WHERE w.id = $1
			  AND (w.privacy = 'public' OR w.author_id = $2 OR ws.user_id IS NOT NULL)`

	weblog := &model.Weblog{}

	err := r.db.QueryRow(query, id,userID).Scan(&weblog.ID, &weblog.Title, &weblog.Content, &weblog.Image, &weblog.AuthorID, &weblog.Privacy, &weblog.CreatedAt)

	if err != nil {
		return nil, err
	}
	return weblog, nil
}

func (r *WeblogRepo) Delete (id int64) error {
	query := `DELETE FROM weblogs WHERE id = $1`

	res, err := r.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("delete weblog: %w", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}