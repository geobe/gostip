package controller

import (
	"github.com/geobe/gostip/go/model"
	scc "github.com/gorilla/securecookie"
	"net/http"
	"log"
)

// filter is called before chaining handlers. Next handler in
// the chain is only called when filter returns true
type filter func(http.ResponseWriter, *http.Request) bool

// a struct to chain several handlers for use with alice
type chainableHandler struct {
	filter filter
	chain  http.Handler
}

// make chainableHandler an http.Handler
func (c chainableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.filter(w, r) {
		c.chain.ServeHTTP(w, r)
	}
}

// SessionChecker filter checks if there is a valid session,
// i.e if someone is logged in
func SessionChecker(h http.Handler) http.Handler {
	c := chainableHandler{chain: h,
		filter: filter(checkSession)}
	return c
}

// here the session check is actually implemented
func checkSession(w http.ResponseWriter, r *http.Request) bool {
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

// RequestLogger uses logRequest function to
// log request info to log output
func RequestLogger(h http.Handler) http.Handler {
	c := chainableHandler{
		filter: filter(logRequest),
		chain:  h,
	}
	return c
}

// logRequest writes relevant information from the request
// to the logging output
func logRequest(w http.ResponseWriter, r *http.Request) bool {
	log.Printf("\t%s: %s%s\n", r.Method, r.Host, r.URL.Path)
	return true
}

// checkAuth is the filter function where the actual authorizing is done
func checkAuth(w http.ResponseWriter, r *http.Request, mask interface{}) bool {
	session, e0 := SessionStore().Get(r, S_DKFAI)
	m, ifaceok := mask.(int)
	role, sessionok := session.Values["role"].(int)
	if e0 != nil || !ifaceok || !sessionok {
		http.Error(w, "error validating role", http.StatusInternalServerError)
		return false
	}
	if role & m == 0 {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return false
	}
	return true
}

// reduce additional function parameter to get a filter
func makeFilter(f func(http.ResponseWriter, *http.Request, interface{}) bool,
		mask interface{}) filter {
	return func(w http.ResponseWriter, r *http.Request) bool {
		return f(w, r, mask)
	}
}

// authorize for anyone who is logged in
func AuthAny(h http.Handler) http.Handler {
	c := chainableHandler{
		filter: makeFilter(checkAuth, model.U_ALL),
		chain:  h,
	}
	return c
}

// authorize for deans office staff for enrolling
func AuthEnrol(h http.Handler) http.Handler {
	c := chainableHandler{
		filter: makeFilter(checkAuth, model.U_ENROL),
		chain:  h,
	}
	return c
}

// authorize for project office staff
func AuthProjectOffice(h http.Handler) http.Handler {
	c := chainableHandler{
		filter: makeFilter(checkAuth, model.U_POFF),
		chain:  h,
	}
	return c
}

// authorize for user administrator
func AuthUserAdmin(h http.Handler) http.Handler {
	c := chainableHandler{
		filter: makeFilter(checkAuth, model.U_UADMIN),
		chain:  h,
	}
	return c
}

// authorize for master administrator
func AuthMasterAdmin(h http.Handler) http.Handler {
	c := chainableHandler{
		filter: makeFilter(checkAuth, model.U_FULLADMIN),
		chain:  h,
	}
	return c
}

//var options  []csrf.Option = []csrf.Option{csrf.Secure(viper.GetBool("csrfsecure"))}
//var protector http.Handler
//var protsem sync.Mutex
//
//func CsrfChecker(h http.Handler) http.Handler {
//	protsem.Lock()
//	defer protsem.Unlock()
//	if protector == nil {
//		key := GetCsrfKey()
//		log.Printf("CsrfChecker created with key %v\n", key)
//		protector = csrf.Protect(key, csrf.Secure(false))(h)
//	} else {
//		log.Println("reusing CsrfChecker")
//	}
//	return protector
//}
//
//func GetCsrfKey() (key []byte) {
//	key = make([]byte, 32)
//	vkey := viper.Get("csrfkey").([]interface{})
//	for i, v := range vkey {
//		if(i >= 32) {
//			break
//		}
//		key[i] = byte(v.(float64))
//	}
//	return
//}
