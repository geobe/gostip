package main

import (
	"log"
	"github.com/geobe/gostip/go/view"
)

// map transports values from go code to templates
type tmodel map[string]interface{}

func main() {
	view.ReadTranslations(view.DEFAULT_TRANSLATIONS_FILE)
	log.Print(view.I18n("_dkfai_welcome", "de"))
	log.Print(view.I18n("_dkfai_welcome", "ru"))
	log.Print(view.I18n("_dkfai_welcome", "en"))
	log.Print(view.I18n("_dkfai_welcome", "ru"))
	values := tmodel{
		"firstname": "Fiete",
		"lastname": "Kall",
		"email": "fiete@kall.net",
	}
	log.Print(view.I18n("_hello_personal", "de", values))
	log.Print(view.I18n("_hello_personal", "ru", values))
	log.Print(view.I18n("_hello_personal", "kg", values))
}
