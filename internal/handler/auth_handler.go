package handler

import (
	"errors"
	"net/http"

	"github.com/NooraPanahi/Weblog-Application.git/internal/service"
	"github.com/NooraPanahi/Weblog-Application.git/internal/session"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService    *service.AuthService
	sessionManager *session.Manager
}

func NewAuthHandler(authService *service.AuthService, sessionManager *session.Manager) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionManager: sessionManager}
}

func (h *AuthHandler) ShowRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func (h *AuthHandler) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.authService.Register(username, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUsernameTaken):
			return c.Render(http.StatusConflict,
				"register.html", map[string]interface{}{
					"Error": "Username is already taken",
				},
			)
		case errors.Is(err, service.ErrInvalidInput):
			return c.Render(http.StatusBadRequest,
				"register.html", map[string]interface{}{
					"Error": "invalid username or password",
				},
			)

		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to register user")
		}
	}

	if err := h.sessionManager.SetUserID(c.Response().Writer, c.Request(), user.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) ShowLogin (c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (h *AuthHandler) Login (c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.authService.Login(username, password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials){ 
			return c.Render(http.StatusUnauthorized,
				"login.html", map[string]string{
					"Error": "Invalid username or password",
				},
			)			
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to login")
		
	}

	if err := h.sessionManager.SetUserID(c.Response().Writer, c.Request(), user.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

	return c.Redirect(http.StatusSeeOther, "/")	
}
