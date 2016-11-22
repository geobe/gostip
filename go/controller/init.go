// Package controller holds all handlers and handler functions
// as well as necessary infrastructure for session management
// and security
package controller

import (
	scc "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var sessionStore = makeStore()

// key for the session store
const S_DKFAI = "DKFAI-App-Session"

// map transports values from go code to templates
type viewmodel map[string]interface{}

// helper function to create a gorilla session store with
// a strong set of keys
func makeStore() sessions.Store {
	return sessions.NewCookieStore(
		scc.GenerateRandomKey(32),
		scc.GenerateRandomKey(32))
}

// accessor for the gorilla session store
func SessionStore() sessions.Store {
	return sessionStore
}
