package repo

import (
	"database/sql"
	"fmt"

	"github.com/NooraPanahi/Weblog-Application.git/internal/model"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (re *CommentRepo) Create(comment *model.Comment) error {
	query := `INSERT INTO comments (weblog_id, user_id, content)
			  VALUES ($1, $2, $3)
			  RETURNING id, created_at`

	err := re.db.QueryRow(query, comment.WeblogID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil

}

func (co *CommentRepo) ListByWeblogID (weblogID int64) ([]*model.Comment, error) {
	query:= `SELECT id, weblog_id, user_id, content, created_at FROM comments
			 WHERE weblog_id = $1 
			 ORDER BY created_at ASC`

	rows, err := co.db.Query(query, weblogID)
	if err != nil {
		return  nil, err
	}

	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		comment := &model.Comment{}
		err := rows.Scan(&comment.ID, &comment.WeblogID, &comment.UserID, &comment.Content, &comment.CreatedAt)

		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}