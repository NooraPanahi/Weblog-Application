package session

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/gorilla/sessions"
)

const (
	sessionName = "weblog_session"
	userIDKey   = "user_id"
)

type Manager struct {
	store *sessions.CookieStore
}

func NewManager(secret string) *Manager {
	return &Manager{store: sessions.NewCookieStore([]byte(secret))}
}

func (m *Manager) SetUserID(c echo.Context, userid int64) error {
	se, err := m.store.Get(c.Request(), sessionName)

	if err != nil {
		return err
	}

	se.Values[userIDKey] = strconv.FormatInt(userid, 10)

	return se.Save(c.Request(), c.Response().Writer)
}

func (m *Manager) GetUserID(c echo.Context) (int64, bool, error) {
	see, err := m.store.Get(c.Request(), sessionName)

	if err != nil {
		return 0, false, err
	}

	value, ok := see.Values[userIDKey]
	if !ok {
		return 0, false, nil
	}

	userIDString, ok := value.(string)

	if !ok {
		return 0, false, nil
	}
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return 0, false, nil
	}

	return userID, true, nil
}

func (m *Manager) Clear(c echo.Context) error {
	see, err := m.store.Get(c.Request(), sessionName)

	if err != nil {
		return err
	}

	see.Options.MaxAge = -1

	return see.Save(c.Request(), c.Response().Writer)
}
