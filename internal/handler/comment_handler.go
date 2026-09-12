package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/NooraPanahi/Weblog-Application.git/internal/service"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct{
	commentService *service.CommentService
}

func NewCommentHandler(commentservice *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentservice}
}

func (h *CommentHandler) Create(c echo.Context) error {
	weblogID, err := strconv.ParseInt(c.Param("id"), 10,64)

	if err != nil || weblogID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid weblog id")
	}

	userID, ok := c.Get("userID").(int64)

	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	content := c.FormValue("content")

	_ , err = h.commentService.Create(weblogID, userID, content)

	if err != nil {
		
		if errors.Is(err, service.ErrInvalidCommentInput){
			return echo.NewHTTPError(http.StatusBadRequest, "invalid comment")
		}

		if errors.Is(err, service.ErrWeblogNotFound){
			return echo.NewHTTPError(http.StatusNotFound, "weblog not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create comment")
	}

	return c.Redirect(http.StatusSeeOther, "/weblog/"+strconv.FormatInt(weblogID,10))

}

func (h *CommentHandler) Delete(c echo.Context) error {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || commentID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid comment id")
	}

	userID, ok := c.Get("userID").(int64)

	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.commentService.Delete(commentID, userID)

	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "comment not found")
		}

		if errors.Is(err, service.ErrorCommentForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, "you cannot delete this comment")
		}

		if errors.Is(err, service.ErrInvalidCommentInput) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid comment")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete comment")
	}

	return c.Redirect(http.StatusSeeOther, c.Request().Referer())
}