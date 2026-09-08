package service

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/NooraPanahi/Weblog-Application.git/internal/model"
	"github.com/NooraPanahi/Weblog-Application.git/internal/repo"
)

var (
	ErrInvalidWeblogInput = errors.New("invalid weblog input")
	ErrInvalidPrivacy     = errors.New("invalid privacy")
	ErrWeblogNotFound     = errors.New("weblog not found")
)

type WeblogService struct {
	weblogRepo *repo.WeblogRepo
}

func NewWeblogService(weblogRepo *repo.WeblogRepo) *WeblogService {
	return &WeblogService{weblogRepo: weblogRepo}
}

func (s *WeblogService) Create(title, content string, image *string, authorID int64, privacy string) (*model.Weblog, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	privacy = strings.TrimSpace(privacy)

	if title == "" || content == "" || authorID <= 0 {
		return nil, ErrInvalidWeblogInput
	}

	if len(title) > 255 {
		return nil, ErrInvalidWeblogInput
	}

	if len(privacy) == 0 {
		privacy = "public"
	}

	if privacy != "public" && privacy != "private" {
		return nil, ErrInvalidPrivacy
	}

	weblog := &model.Weblog{Title: title, Content: content, Image: image, AuthorID: authorID, Privacy: privacy}

	if err := s.weblogRepo.Create(weblog); err != nil {
		return nil, err
	}

	return weblog, nil
}

func (s *WeblogService) ListVisible(userID int64) ([]*model.Weblog, error) {
	if userID <= 0 {
		return nil, ErrInvalidWeblogInput
	}

	weblogs, err := s.weblogRepo.ListVisible(userID)

	if err != nil {
		return nil, err
	}
	return weblogs, nil
}

func (s *WeblogService) GetVisibleWeblogByID(weblogID, userID int64) (*model.Weblog, error) {
	if weblogID <= 0 || userID <= 0 {
		return nil, ErrInvalidWeblogInput
	}

	weblog, err := s.weblogRepo.FindVisibleByID(weblogID, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWeblogNotFound
		}
		return nil, err
	}
	return weblog, nil
}
