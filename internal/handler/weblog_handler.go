package handler

import (
	"errors"
	"net/http"

	"github.com/NooraPanahi/Weblog-Application.git/internal/service"
	"github.com/labstack/echo/v4"
)

type WeblogHandler struct {
	weblogService *service.WeblogService
}

func NewWeblogHandler(weblogservice *service.WeblogService) *WeblogHandler {
	return &WeblogHandler{weblogService: weblogservice}
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

	var image *string
	imageValue := c.FormValue("image")

	if imageValue != "" {
		image = &imageValue
	}

	_, err := h.weblogService.Create(title, content, image, userID, privacy)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Render(http.StatusBadRequest, "create.html", map[string]string{"Error": "Invalid weblog data"})

		case errors.Is(err, service.ErrInvalidPrivacy):
			return c.Render(http.StatusBadRequest, "create.html", map[string]string{"Error": "Invalid privacy"})

		default: return echo.NewHTTPError(http.StatusInternalServerError, "failed to create weblog")
		}
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
