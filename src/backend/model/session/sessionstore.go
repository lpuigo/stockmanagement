package session

import (
	"github.com/gorilla/sessions"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/user"
	"net/http"
)

const (
	sessionName string = "Batec-Session"
	userId      string = "userid"

	sessionMaxAge int = 86400 * 7
)

type SessionStore struct {
	*sessions.CookieStore
	SessionName string
}

func NewSessionStore(secretkey string) *SessionStore {
	ss := &SessionStore{}
	ss.CookieStore = sessions.NewCookieStore([]byte(secretkey))
	ss.SessionName = sessionName
	return ss
}

// CheckUser checks for request Session and return connected User Id (-1 if no user properly connected)
func (ss *SessionStore) CheckUser(r *http.Request) int {
	session, _ := ss.Get(r, ss.SessionName)
	if session.IsNew {
		return -1
	}
	user, ok := session.Values[userId].(int)
	if !ok {
		return -1
	}
	return user
}

func (ss *SessionStore) RefreshSessionCookie(w http.ResponseWriter, r *http.Request) error {
	session, err := ss.Get(r, ss.SessionName)
	if err != nil {
		return err
	}
	// update session MaxAge
	session.Options.MaxAge = sessionMaxAge
	// Save it before we write to the response/return from the handler.
	return session.Save(r, w)
}

func (ss *SessionStore) AddSessionCookie(u *user.UserRecord, w http.ResponseWriter, r *http.Request) error {
	session, err := ss.Get(r, ss.SessionName)
	if err != nil {
		return err
	}
	// Set some session values.
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		HttpOnly: true,
	}
	session.Values[userId] = u.Id
	// Save it before we write to the response/return from the handler.
	return session.Save(r, w)
}

func (ss *SessionStore) RemoveSessionCookie(w http.ResponseWriter, r *http.Request) error {
	session, err := ss.Get(r, ss.SessionName)
	if err != nil {
		return err
	}
	//session.Options = &sessions.Options{
	//	Path:     "/",
	//	MaxAge:   -1,
	//	HttpOnly: true,
	//}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
