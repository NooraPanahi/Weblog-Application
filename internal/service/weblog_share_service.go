package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/NooraPanahi/Weblog-Application.git/internal/repo"
)

var (
	ErrInvalidShareInput = errors.New("invalid share input")
	ErrShareNotWeblogOwner    = errors.New("user is not weblog owner")
	ErrShareUserNotFound = errors.New("share user not found")
	ErrAlreadyShared     = errors.New("weblog already shared with user")
)

type WeblogShareService struct {
	shareRepo  *repo.WeblogShareRepo
	userRepo   *repo.UserRepo
	weblogRepo *repo.WeblogRepo
}

func NewWeblogShareService(shareR *repo.WeblogShareRepo, userR *repo.UserRepo, weblogR *repo.WeblogRepo) *WeblogShareService {
	return &WeblogShareService{
		shareRepo:  shareR,
		userRepo:   userR,
		weblogRepo: weblogR,
	}
}

func (s *WeblogShareService) Share (weblogID, ownerID int64, username string)error {
	username = strings.TrimSpace(username)

	if weblogID <= 0 || ownerID <= 0 || username == "" {
		return ErrInvalidShareInput
	}

	weblog, err := s.weblogRepo.FindWeblogByID(weblogID)

	if err != nil {
		return fmt.Errorf("find weblog: %w", err)
	}

	if weblog.AuthorID != ownerID {
		return ErrShareNotWeblogOwner
	}	

	user, err := s.userRepo.FindUserByUsername(username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrShareUserNotFound
		}
		return fmt.Errorf("find share user: %w", err)
	}

	if user.ID == ownerID {
		return ErrInvalidShareInput
	}

	exists, err := s.shareRepo.Exists(weblogID, user.ID)

	if err != nil {
		return fmt.Errorf("check existing share:%w", err)
	}

	if exists {
		return ErrAlreadyShared
	}

	if err := s.shareRepo.Create(weblogID, user.ID); err != nil {
		return fmt.Errorf("share weblog: %w", err)
	}
	return nil
}