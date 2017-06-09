package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	scc "github.com/gorilla/securecookie"
	"net/http"
	"github.com/justinas/nosurf"
	"log"
	"html"
)

// func HandleLogin handles requests to login.
// For GET requests, the login form gets displayed.
// POST requests are checked for a valid username/password.
// If check succeeds, a user session is initiated and username
// and user information are stored in the session store.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f := nosurf.Token(r)
		values := viewmodel{
			"csrftoken": f,
		}
		view.Views().ExecuteTemplate(w, "login", values)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		login := html.EscapeString(r.PostFormValue("login"))
		passwd := html.EscapeString(r.PostFormValue("password"))
		db := model.Db()
		var user model.User
		db.First(&user, "login = ?", login)
		if user.ID > 0 && user.ValidatePw(passwd) {
			err := createUserSession(user, w, r)
			if err == nil {
				log.Printf("User %s logged in\n", login)
				http.Redirect(w, r, "/work", http.StatusFound)
			} else {
				http.Error(w, "Session error: " + err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			view.Views().ExecuteTemplate(w, "login",
				map[string]string{"failure": "Failure, wrong login or password"})
		}
	}
}

// helper function to create a new session and store user related info in the session store.
func createUserSession(user model.User, w http.ResponseWriter, r *http.Request) (err error) {
	// let's start or retrieve a session
	session, err := SessionStore().Get(r, S_DKFAI)
	if err != nil {
		if err.(scc.Error).IsDecode() {
			// no valid session, reset err and overwrite session
			err = nil
			session, _ = SessionStore().New(r, S_DKFAI)
		} else {
			http.Error(w, "Get error: " + err.Error(),
				http.StatusInternalServerError)
			return
		}
	}
	session.Values["login"] = user.Login
	session.Values["fullname"] = user.Fullname
	session.Values["role"] = user.Role
	session.Values["userid"] = user.ID
	session.Values["language"] = view.PreferedLanguages(r) [0]
	session.Options.MaxAge = 60 * 20 // limit session to 20 min
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Save error: " + err.Error(), http.StatusInternalServerError)
	}
	return
}

// helper function to get user information for current session
// from session store
func getUserFromSession(r *http.Request) (model.User, error) {
	session, err := SessionStore().Get(r, S_DKFAI)
	user := model.User{}
	uid := session.Values["userid"]
	if err == nil && uid != nil {
		user.ID = session.Values["userid"].(uint)
		user.Login = session.Values["login"].(string)
		user.Fullname = session.Values["fullname"].(string)
		user.Role = session.Values["role"].(int)
	}
	return user, err
}

// func HandleLogout serves requests to /logout. It cancels the current
// user session and redirects to login page.
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := SessionStore().Get(r, S_DKFAI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// delete the session
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
	return
}
