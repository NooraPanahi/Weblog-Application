package service

import (
	"errors"
	"strings"

	"github.com/NooraPanahi/Weblog-Application.git/internal/model"
	"github.com/NooraPanahi/Weblog-Application.git/internal/repo"
)

var (
	ErrInvalidWeblogInput = errors.New("invalid weblog input")
	ErrInvalidPrivacy     = errors.New("invalid privacy")
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

