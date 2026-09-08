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
}

func NewWeblogShareHandler(shareS *service.WeblogShareService) *WeblogShareHandler {
	return &WeblogShareHandler{shareService: shareS}
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
			return echo.NewHTTPError(http.StatusBadRequest, "invalid share input")

		case errors.Is(err, service.ErrNotWeblogOwner):
			return echo.NewHTTPError(http.StatusForbidden, "only the weblog owner can share it")

		case errors.Is(err, service.ErrShareUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "user not found")

		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to share weblog")
		}
	}

	return c.Redirect(http.StatusSeeOther,"/weblog/"+strconv.FormatInt(weblogID, 10))
}
