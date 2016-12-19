// Package controller holds all handlers and handler functions
// as well as necessary infrastructure for session management
// and security
package controller

import (
	scc "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"encoding/gob"
	"github.com/geobe/gostip/go/model"
)

var sessionStore = makeStore()

// keys for the session store
const S_DKFAI = "DKFAI-App-Session"
const S_APPLICANT = "GOSTIP_APPLICANT"
const S_USER = "GOSTIP_USER"

// map transports values from go code to templates
type viewmodel map[string]interface{}

// helper function to create a gorilla session store with
// a strong set of keys
func makeStore() sessions.Store {
	store := sessions.NewCookieStore(
		scc.GenerateRandomKey(32),
		scc.GenerateRandomKey(32))
	registerTypes()
	return store
}

// accessor for the gorilla session store
func SessionStore() sessions.Store {
	return sessionStore
}

func registerTypes() {
	gob.Register(model.Applicant{})
	gob.Register(model.ApplicantData{})
	gob.Register(model.Oblast{})
	gob.Register(model.User{})
}