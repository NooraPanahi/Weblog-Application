package repo

import (
	"database/sql"
	"fmt"

	"github.com/NooraPanahi/Weblog-Application.git/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (us *UserRepo) Create(user *model.User) error {
	query := `INSERT INTO users (username, password)
			  VALUES ($1, $2)
			  RETURNING id, created_at`

	err := us.db.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (us *UserRepo) FindUserByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, password, created_at
			  FROM users WHERE username = $1`

	var user model.User

	err := us.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &user, nil
}

func (us *UserRepo) FindUserByID(id int64) (*model.User, error) {
	query := `SELECT id, username, password, created_at
			  FROM users WHERE id = $1`

	var user model.User

	err := us.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}