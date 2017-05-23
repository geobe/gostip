package controller

import (
	"net/http"
	"net/url"
	"strings"
	"log"
	"regexp"
)


// RedirectHTTP is an HTTP handler (suitable for use with http.HandleFunc)
// that responds to all requests by redirecting to the same URL served over HTTPS.
// It should only be invoked for requests received over HTTP.
func RedirectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil || r.Host == "" {
		http.Error(w, "not found", 404)
	}

	var u *url.URL
	u = r.URL
	host := r.Host
	u.Host = strings.Split(host, ":")[0]
	u.Scheme = "https"
	log.Printf("redirect to u.host  %s -> %s\n", r.Host, u.String())
	http.Redirect(w, r, u.String(), 302)
}

// interne Aufrufe vom gleichen lokalen Netz mit Mux annehmen, sonst redirect auf HTTPS
type HandlerSwitch struct {
	Mux      http.Handler
	Redirect http.Handler
}

// nicht richtig für 172.16.0.0/12
var matcher = regexp.MustCompile("(192\\.168.*)|(localhost)|(10\\..*)|(172\\..*)")

// Handler function, die Internetzugriffe auf den Redirect Handler umleitet.
// Lokale Zugriffe werden direkt von MUX geroutet
func (h *HandlerSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	local := matcher.MatchString(host)
	if local {
		h.Mux.ServeHTTP(w, r)
	} else {
		h.Redirect.ServeHTTP(w, r)
	}
}
