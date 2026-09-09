package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/NooraPanahi/Weblog-Application.git/internal/service"
	"github.com/NooraPanahi/Weblog-Application.git/internal/upload"
	"github.com/labstack/echo/v4"
)

type WeblogHandler struct {
	weblogService  *service.WeblogService
	commentService *service.CommentService
}

func NewWeblogHandler(weblogservice *service.WeblogService, commentservice *service.CommentService) *WeblogHandler {
	return &WeblogHandler{weblogService: weblogservice, commentService: commentservice}
}

func (h *WeblogHandler) ShowCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "create.html", nil)
}

func (h *WeblogHandler) Create(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	privacy := c.FormValue("privacy")

	userID, ok := c.Get("userID").(int64)

	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	file, err := c.FormFile("image")

	if err != nil {
		file = nil
	}
	var image *string

	if file != nil {
		imagePath, err := upload.SaveImage(file)

		if err != nil {
			return c.Render(http.StatusBadRequest, "create.html", map[string]string {"Error": err.Error()})
		}
		image = &imagePath
	}
	_, err = h.weblogService.Create(title, content, image, userID, privacy)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Render(http.StatusBadRequest, "create.html", map[string]string{"Error": "Invalid weblog data"})

		case errors.Is(err, service.ErrInvalidPrivacy):
			return c.Render(http.StatusBadRequest, "create.html", map[string]string{"Error": "Invalid privacy"})

		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create weblog")
		}
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *WeblogHandler) Detail(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid weblog id")
	}

	userID, ok := c.Get("userID").(int64)

	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	weblog, err := h.weblogService.GetVisibleWeblogByID(id, userID)

	if err != nil {
		if errors.Is(err, service.ErrWeblogNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Weblog not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load weblog")
	}

	comments, err := h.commentService.ListByWeblogID(id, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to  load comments")
	}
	data := map[string]interface{}{
		"Weblog":   weblog,
		"Comments": comments,
		"UserID":   userID,
	}

	return c.Render(http.StatusOK, "detail.html", data)
}

func (h *WeblogHandler) Delete(c echo.Context) error {
	weblogID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || weblogID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid weblog id")
	}

	userID, ok := c.Get("userID").(int64)
	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.weblogService.Delete(weblogID, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidWeblogInput):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid weblog")

		case errors.Is(err, service.ErrWeblogNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "weblog not found")

		case errors.Is(err, service.ErrNotWeblogOwner):
			return echo.NewHTTPError(http.StatusForbidden,"only the weblog owner can delete it")

		default:
			return echo.NewHTTPError(http.StatusInternalServerError,"failed to delete weblog")
		}
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
