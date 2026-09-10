package handler

import (
	"net/http"

	"github.com/NooraPanahi/Weblog-Application.git/internal/service"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
	weblogService *service.WeblogService
}

func NewHomeHandler(weblogservice *service.WeblogService) *HomeHandler {
	return &HomeHandler{weblogService: weblogservice}
}

func (h *HomeHandler) Home(c echo.Context) error {
	userID, ok := c.Get("userID").(int64)

	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	MyWeblogs, err := h.weblogService.ListMyWeblogs(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load weblogs")
	}

	SharedWeblogs, err := h.weblogService.ListSharedWeblogs(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load shared weblogs")
	}

	publicWeblogs, err := h.weblogService.ListPublicWeblogs(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load public weblogs")
	}

	data := map[string]interface{}{
		"MyWeblogs":     MyWeblogs,
		"SharedWeblogs": SharedWeblogs,
		"PublicWeblogs": publicWeblogs,
		"Username":      c.Get("username"),
	}

	return c.Render(http.StatusOK, "home.html", data)
}
