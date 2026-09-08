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