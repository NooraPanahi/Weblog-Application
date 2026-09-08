package repo

import (
	"database/sql"
	"fmt"
)

type WeblogShareRepo struct {
	db *sql.DB
}

func NewWeblogShareRepo(db *sql.DB) *WeblogShareRepo {
	return &WeblogShareRepo{db: db}
}

func (r *WeblogShareRepo) Create (weblogID, userID int64) error {
	query := `INSERT INTO weblog_shares (weblog_id, user_id)
			  VALUES ($1, $2)`
	_, err := r.db.Exec(query, weblogID, userID)

	if err != nil {
		return fmt.Errorf("create weblog share: %w", err)
	}

	return nil
}
