package service

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/NooraPanahi/Weblog-Application.git/internal/model"
	"github.com/NooraPanahi/Weblog-Application.git/internal/repo"
	"golang.org/x/crypto/bcrypt"
)


var ( 
	ErrUsernameTaken      = errors.New("username already taken") 
	ErrUserNotFound       = errors.New("user not found") 
	ErrInvalidInput       = errors.New("invalid input") 
	ErrInvalidCredentials = errors.New("invalid username or password") 
)


type AuthService struct {
	userRepo *repo.UserRepo
}

func NewAuthService(userRepo *repo.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register (username, password string) (*model.User , error) {
	username = strings.TrimSpace(username)

	if username == "" || password == "" {
		return nil, ErrInvalidInput
	}

	if len(password) < 6 {
		return nil, ErrInvalidInput
	}
	_,err := s.userRepo.FindUserByUsername(username)
	if err == nil {
		return nil, ErrUsernameTaken
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User {
		Username: username, PasswordHash: string(passwordHash),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}


func (s *AuthService) Login(username, password string) (*model.User, error) {
	username = strings.TrimSpace(username)

	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}


	user, err := s.userRepo.FindUserByUsername(username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *AuthService) GetUserByID (id int64) (*model.User, error) {
	user, err := s.userRepo.FindUserByID(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	} 

	return user, nil
}