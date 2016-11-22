package controller

import (
	"github.com/geobe/gostip/go/model"
	scc "github.com/gorilla/securecookie"
	"net/http"
)

// chainfunc is called before chaining handlers. Next handler in
// the chain is only called when chainfunc returns true
type chainfunc func(http.ResponseWriter, *http.Request, interface{}) bool

// a struct to chain several handlers for use with alice
type chainableHandler struct {
	filter chainfunc
	chain  http.Handler
}

// make chainableHandler an http.Handler
func (c chainableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.filter(w, r, nil) {
		c.chain.ServeHTTP(w, r)
	}
}

// SessionChecker filter checks if there is a valid session,
// i.e if someone is logged in
func SessionChecker(h http.Handler) http.Handler {
	c := chainableHandler{chain: h,
		filter: chainfunc(checkSession)}
	return c
}

// here the session check is actually implemented
func checkSession(w http.ResponseWriter, r *http.Request, ignore interface{}) bool {
	session, err := SessionStore().Get(r, S_DKFAI)
	if err != nil {
		if err.(scc.Error).IsDecode() {
			// recover from an old hanging session going to login
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return false
	}
	if session.IsNew {
		// no session there, goto login
		http.Redirect(w, r, "/login", http.StatusFound)
		return false
	}
	return true
}

// extend chainableHandler for authorisation
type authHandler struct {
	chainableHandler
	authMask int
}

// make chainableHandler an http.Handler
func (c authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.filter(w, r, c.authMask) {
		c.chain.ServeHTTP(w, r)
	}
}

// create an authorization filter handler
func authorizer(mask int, h http.Handler) http.Handler {
	return authHandler{
		chainableHandler{chainfunc(checkAuth), h},
		mask,
	}
}

// here the actual authorizing is done
func checkAuth(w http.ResponseWriter, r *http.Request, mask interface{}) bool {
	session, e0 := SessionStore().Get(r, S_DKFAI)
	m, e1 := mask.(int)
	role, e2 := session.Values["role"].(int)
	if e0 != nil || !e1 || !e2 {
		http.Error(w, "error validating role", http.StatusInternalServerError)
		return false
	}
	if role&m == 0 {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return false
	}
	return true
}

// authorize for anyone who is logged in
func AuthAny(h http.Handler) http.Handler {
	return authorizer(model.U_ALL, h)
}

// authorize for deans office staff for enrolling
func AuthEnrol(h http.Handler) http.Handler {
	return authorizer(model.U_ENROL, h)
}

// authorize for project office staff
func AuthProjectOffice(h http.Handler) http.Handler {
	return authorizer(model.U_POFF, h)
}

// authorize for user administrator
func AuthUserAdmin(h http.Handler) http.Handler {
	return authorizer(model.U_UADMIN, h)
}

// authorize for master administrator
func AuthMasterAdmin(h http.Handler) http.Handler {
	return authorizer(model.U_FULLADMIN, h)
}
