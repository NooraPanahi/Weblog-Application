package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/NooraPanahi/Weblog-Application.git/internal/service"
	"github.com/labstack/echo/v4"
)

type WeblogShareHandler struct {
	shareService *service.WeblogShareService
	weblogService  *service.WeblogService
	commentService *service.CommentService
}

func NewWeblogShareHandler(shareS *service.WeblogShareService, weblogS *service.WeblogService, commentS *service.CommentService) *WeblogShareHandler {
	return &WeblogShareHandler{
		shareService:   shareS,
		weblogService:  weblogS,
		commentService: commentS,
	}
}

func (h *WeblogShareHandler) Share(c echo.Context) error {
	weblogID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || weblogID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid weblog id")
	}
	userID, ok := c.Get("userID").(int64)

	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	username := c.FormValue("username")

	err = h.shareService.Share(weblogID, userID, username)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidShareInput):
			return h.renderDetailWithError(c, weblogID, userID, "Invalid share input")
			
		case errors.Is(err, service.ErrNotWeblogOwner):
			return echo.NewHTTPError(http.StatusForbidden, "only the weblog owner can share it")

		case errors.Is(err, service.ErrShareUserNotFound):
			return h.renderDetailWithError(c, weblogID, userID, "User not found")
			
		case errors.Is(err, service.ErrAlreadyShared):
			return h.renderDetailWithError(c, weblogID, userID, "This weblog is already shared with this user")
			
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to share weblog")
		}
	}

	return c.Redirect(http.StatusSeeOther,"/weblog/"+strconv.FormatInt(weblogID, 10))
}

func (h *WeblogShareHandler) renderDetailWithError(c echo.Context, weblogID int64, userID int64, message string) error {

	weblog, err := h.weblogService.GetVisibleWeblogByID(weblogID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "weblog not found")
	}

	comments, err := h.commentService.ListByWeblogID(weblogID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load comments")
	}

	return c.Render(http.StatusBadRequest, "detail.html", map[string]interface{}{
		"Weblog":   weblog,
		"Comments": comments,
		"UserID":   userID,
		"Error":    message,
	})
}