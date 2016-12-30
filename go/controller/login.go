package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	scc "github.com/gorilla/securecookie"
	"net/http"
	"log"
	"github.com/justinas/nosurf"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f := nosurf.Token(r)
		log.Printf("token: %v\n", f)
		values := viewmodel{
			"csrftoken": f,
		}
		view.Views().ExecuteTemplate(w, "login", values)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		login := r.PostFormValue("login")
		passwd := r.PostFormValue("password")
		db := model.Db()
		var user model.User
		db.First(&user, "login = ?", login)
		if user.ID > 0 && user.ValidatePw(passwd) {
			// let's start or retrieve a session
			session, err := SessionStore().Get(r, S_DKFAI)
			if err != nil {
				if err.(scc.Error).IsDecode() {
					session, _ = SessionStore().New(r, S_DKFAI)
				} else {
					http.Error(w, "Get error: "+err.Error(),
						http.StatusInternalServerError)
					return
				}
			}
			session.Values["login"] = user.Login
			session.Values["fullname"] = user.Fullname
			session.Values["role"] = user.Role
			session.Options.MaxAge = 60 * 20 // limit session to 20 min
			if err := session.Save(r, w); err != nil {
				http.Error(w, "Save error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/work", http.StatusFound)
			return
		} else {
			//http.Error(w, "Login Failure", http.StatusForbidden)
			view.Views().ExecuteTemplate(w, "login",
				map[string]string{"failure": "Failure, wrong login or password"})
		}
	}
}

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
