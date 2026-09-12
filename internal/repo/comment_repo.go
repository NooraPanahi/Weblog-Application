package repo

import (
	"database/sql"
	"errors"
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

func (co *CommentRepo) ListByWeblogID (weblogID int64) ([]*model.CommentView, error) {
	query:= `SELECT c.id, c.weblog_id, c.user_id, u.username ,  c.content, c.created_at FROM comments c
			 JOIN users u
			 	ON u.id = c.user_id
			 WHERE c.weblog_id = $1
			 ORDER BY c.created_at ASC`

	rows, err := co.db.Query(query, weblogID)
	if err != nil {
		return  nil, err
	}

	defer rows.Close()

	var comments []*model.CommentView

	for rows.Next() {
		comment := &model.CommentView{}
		err := rows.Scan(&comment.ID, &comment.WeblogID, &comment.UserID,&comment.Username, &comment.Content, &comment.CreatedAt)

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

func (co *CommentRepo) GetByID (commentID int64) (*model.Comment, error) {
	query := `SELECT id, weblog_id, user_id, content, created_at
			  FROM comments WHERE id = $1`

	comment := &model.Comment{}
	err := co.db.QueryRow(query, commentID).Scan(&comment.ID, &comment.WeblogID,
		&comment.UserID, &comment.Content,&comment.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get comment: %w", err)
	}
	return comment, nil
}

func (co *CommentRepo) Delete(commentID int64) error {
	query := `DELETE FROM comments WHERE id = $1`

	_, err := co.db.Exec(query, commentID)

	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	return nil
}