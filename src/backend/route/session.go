package route

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	mgr "github.com/lpuig/batec/stockmanagement/src/backend/manager"
	"net/http"
)

type authentUser struct {
	Name        string
	Permissions map[string]bool
}

func newAuthentUser() authentUser {
	return authentUser{
		Name:        "",
		Permissions: make(map[string]bool),
	}
}

func (au *authentUser) SetFrom(mgr *mgr.Manager) {
	au.Name = mgr.CurrentUser.Name
	au.Permissions = mgr.CurrentUser.Permissions
}

// GetUser checks for session cookie, and returns pertaining user
//
// If user is not authenticated, the session is removed
func GetUser(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetUser")
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	// check for session cookie
	if !mgr.CheckSessionUser(r) {
		// user cookie not found or improper, remove it first
		err := mgr.SessionStore.RemoveSessionCookie(w, r)
		if err != nil {
			AddError(w, logmsg, "could not remove session info", http.StatusInternalServerError)
			return
		}
		AddError(w, logmsg, "user not authorized", http.StatusUnauthorized)
		return
	}

	// found a correct one, set user
	logmsg.AddUser(mgr.CurrentUser.Name)
	user := newAuthentUser()
	user.SetFrom(mgr)

	// refresh session cookie
	err := mgr.SessionStore.RefreshSessionCookie(w, r)
	if err != nil {
		AddError(w, logmsg, "could not refresh session cookie", http.StatusInternalServerError)
		return
	}

	// write response
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		AddError(w, logmsg, "could not encode authent user", http.StatusInternalServerError)
		return
	}
	logmsg.Info = "authenticated"
	logmsg.Response = http.StatusOK
}

func Login(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("Login")
	defer logmsg.Log()

	if r.ParseMultipartForm(1024) != nil {
		AddError(w, logmsg, "user info missing", http.StatusBadRequest)
		return
	}

	getValue := func(key string) (string, bool) {
		info, found := r.MultipartForm.Value[key]
		if !found {
			return "", false
		}
		return info[0], true
	}

	user, hasUser := getValue("user")
	pwd, hasPwd := getValue("pwd")
	if !(hasUser && hasPwd) {
		AddError(w, logmsg, "user/password info missing", http.StatusBadRequest)
		return
	}
	logmsg.AddUser(user)

	//TODO Improve Login/pwd and authorization here
	u := mgr.Users.GetByName(user)
	if u == nil {
		AddError(w, logmsg, "user not authorized", http.StatusUnauthorized)
		return
	}
	if u.Password != pwd {
		AddError(w, logmsg, "wrong password", http.StatusUnauthorized)
		return
	}

	if err := mgr.SessionStore.AddSessionCookie(u, w, r); err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
	logmsg.Info = "logged in"
}

func Logout(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("Logout").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	err := mgr.SessionStore.RemoveSessionCookie(w, r)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Info = "logged out"
	logmsg.Response = http.StatusOK
}
