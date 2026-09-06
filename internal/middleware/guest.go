package middleware

import (
	"github.com/NooraPanahi/Weblog-Application.git/internal/session"
	"github.com/labstack/echo/v4"
)

func RequireGuest(sessionManager *session.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, ok, err := sessionManager.GetUserID(c.Request())

			if err != nil {
				return echo.NewHTTPError(500, "failed to read session")
			}

			if ok {
				return c.Redirect(302, "/")
			}
			return next(c)
		}

	}
}
