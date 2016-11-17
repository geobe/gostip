// Package controller holds all handlers and handler functions
// as well as necessary infrastructure for session management
// and security
package controller

import (
	scc "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var sessionStore = makeStore()

const S_DKFAI = "DKFAI-App-Session"

type viewmodel map[string]interface{}
type appvm struct {
	number      uint `gorm:"AUTO_INCREMENT"`
	lastname    string
	firstname   string
	fathersname string
	phone       string
	email       string
	home        string
	school      string
	schoolok    bool
	district    string
	districtok  bool
	ortsum      int16
	ortmath     int16
	ortphys     int16
	ortok       bool
}

func makeStore() sessions.Store {
	return sessions.NewCookieStore(
		scc.GenerateRandomKey(32),
		scc.GenerateRandomKey(32))
}

func SessionStore() sessions.Store {
	return sessionStore
}
