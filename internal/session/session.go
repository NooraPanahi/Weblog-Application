package session

import (
	"net/http"
	"strconv"

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
	store := sessions.NewCookieStore([]byte(secret))

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	return &Manager{store: store}
}

func (m *Manager) SetUserID(w http.ResponseWriter, r *http.Request, userid int64) error {
	se, err := m.store.Get(r, sessionName)

	if err != nil {
		return err
	}

	se.Values[userIDKey] = strconv.FormatInt(userid, 10)

	return se.Save(r, w)
}

func (m *Manager) GetUserID(r *http.Request) (int64, bool, error) {
	see, err := m.store.Get(r, sessionName)

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

func (m *Manager) Clear(w http.ResponseWriter, r *http.Request) error {
	see, err := m.store.Get(r, sessionName)

	if err != nil {
		return err
	}

	see.Options.MaxAge = -1

	return see.Save(r, w)
}
