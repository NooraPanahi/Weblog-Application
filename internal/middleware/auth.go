package middleware

import (
	"net/http"

	"github.com/NooraPanahi/Weblog-Application.git/internal/repo"
	"github.com/NooraPanahi/Weblog-Application.git/internal/session"
	"github.com/labstack/echo/v4"
)

func RequireAuth(sessionManager *session.Manager, userRepo *repo.UserRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok, err := sessionManager.GetUserID(c.Request())

			if err != nil {
				return echo.NewHTTPError(500, "failed to read session")
			}

			if !ok {
				return c.Redirect(302, "/login")
			}

			user, err := userRepo.FindUserByID(userID)

			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/login")
			}
			c.Set("userID", user.ID)
			c.Set("username", user.Username)
			return next(c)

		}

	}
}
