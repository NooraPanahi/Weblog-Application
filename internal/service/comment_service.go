package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NooraPanahi/Weblog-Application.git/internal/model"
	"github.com/NooraPanahi/Weblog-Application.git/internal/repo"
)

var (
	ErrInvalidCommentInput = errors.New("Invalid comment input")
	ErrCommentNotFound     = errors.New("comment not found")
	ErrorCommentForbidden  = errors.New("comment forbidden")
)

type CommentService struct {
	commentRepo   *repo.CommentRepo
	WeblogService *WeblogService
}

func NewCommentService(commentRepo *repo.CommentRepo, weblogservice *WeblogService) *CommentService {
	return &CommentService{commentRepo: commentRepo, WeblogService: weblogservice}
}

func (s *CommentService) Create(weblogID, userID int64, content string) (*model.Comment, error) {
	content = strings.TrimSpace(content)

	if weblogID <= 0 || userID <= 0 || content == "" {
		return nil, ErrInvalidCommentInput
	}
	_, err := s.WeblogService.GetVisibleWeblogByID(weblogID, userID)

	if err != nil {
		return nil, err
	}

	comment := &model.Comment{WeblogID: weblogID, UserID: userID, Content: content}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) ListByWeblogID(weblogID, userID int64) ([]*model.CommentView, error) {
	if weblogID <= 0 || userID <= 0 {
		return nil, ErrInvalidCommentInput
	}

	_, err := s.WeblogService.GetVisibleWeblogByID(weblogID, userID)

	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.ListByWeblogID(weblogID)

	if err != nil {
		return nil, err
	}

	iranLocation, err := time.LoadLocation("Asia/Tehran")

	if err != nil {
		return nil, fmt.Errorf("load timezoan: %w", err)
	}
	for i := range comments {
		comments[i].CreatedAtFormatted = comments[i].CreatedAt.In(iranLocation).Format("Jan 2, 2006 - 15:04")
	}
	return comments, err
}

func (s *CommentService) Delete(commentID, userID int64) error {
	if commentID <= 0 || userID <= 0 {
		return ErrInvalidCommentInput
	}
	comment, err := s.commentRepo.GetByID(commentID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCommentNotFound
		}
		return fmt.Errorf("get comment for delete: %w", err)
	}

	if comment.UserID != userID {
		return ErrorCommentForbidden
	}
	if err := s.commentRepo.Delete(commentID); err != nil {
		return err
	}
	return nil
}
