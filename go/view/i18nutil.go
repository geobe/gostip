package view

import (
	"github.com/geobe/gostip/go/model"
	"os"
	"log"
	"encoding/csv"
	"regexp"
	"net/http"
	"strings"
	"strconv"
	"html/template"
	"bytes"
)

const DEFAULT_TRANSLATIONS_FILE = "allInOne.csv"
const DEFAULT_LANGUAGE = "en"

var translations map[string]map[string]string
var trtemplates map[string]map[string]*template.Template

// is translation a template?
var t4template = regexp.MustCompile(".*\\{\\{.+}}.*")

func ReadTranslations(file string) {
	translations, trtemplates = readCsv(file)
}

func readCsv(file string) (map[string]map[string]string, map[string]map[string]*template.Template) {
	var translations = make(map[string]map[string]string)
	var trtemplates = make(map[string]map[string]*template.Template)
	// Create a function map with dict and refer function
	funcs := map[string]interface{}{
		"safeatt": SafeAtt,
	}

	configfile := model.Base + "/config/" + file
	if rd, err := os.OpenFile(configfile, os.O_RDONLY, 0666); err == nil {
		csvReader := csv.NewReader(rd)
		if all, er2 := csvReader.ReadAll(); er2 == nil {
			head := all[0]
			for _, lang := range head[1:] {
				translations[lang] = make(map[string]string)
				trtemplates[lang] = make(map[string]*template.Template)
			}
			body := all[1:]
			for _, line := range body {
				key := line[0]
				for i, val := range line[1:] {
					lngkey := head[i+1]
					translations[lngkey][key] = val
					if t4template.MatchString(val) {
						t := template.New(key).Funcs(funcs)
						tmpl, err := t.Parse(val)
						if err != nil { panic(err) }
						trtemplates[lngkey][key] = tmpl
					}
				}
			}
		} else {
			log.Printf ("read error %x", er2)
		}
		if translations[DEFAULT_LANGUAGE] != nil {
			translations["default"] = translations[DEFAULT_LANGUAGE]
			trtemplates["default"] = trtemplates[DEFAULT_LANGUAGE]
		} else {
			for k, _ := range translations {
				translations["default"] = translations[k]
			}
			for k, _ := range trtemplates {
				trtemplates["default"] = trtemplates[k]
			}
		}
	}else {
		log.Printf("open error %x", err)
	}
	return translations, trtemplates
}

func GetTranslations(lang string) map[string]string {
	if translations == nil {
		ReadTranslations(DEFAULT_TRANSLATIONS_FILE)
	}
	if translations[lang] != nil {
		return translations[lang]
	} else {
		return  translations["default"]
	}
}

func GetTranslation(key, lang string) string {
	tr := GetTranslations(lang) [key]
	if tr == "" {
		tr = GetTranslations("default") [key]
	}
	if tr == "" {
		tr = key
	}
	return tr
}

func OblastsI18n(lang string) []model.Oblast {
	o := model.Oblasts()
	ointl := make([]model.Oblast, len(o))
	for i, v := range o {
		ointl[i] = model.Oblast{ v.ID, I18n(v.Name, lang) }
	}
	return ointl
}

func LanguagesI18n(lang string) map[model.Lang]string {
	langs := make(map[model.Lang]string)
	for k, v := range model.InitialLanguages {
		langs[k] = I18n(v, lang)
	}
	return langs
}

func GetTrtemplates(lang string) map[string]*template.Template {
	if trtemplates == nil {
		ReadTranslations(DEFAULT_TRANSLATIONS_FILE)
	}
	if trtemplates[lang] != nil {
		return trtemplates[lang]
	} else {
		return  trtemplates["default"]
	}
}

func GetTrtemplate(key, lang string) *template.Template {
	tr := GetTrtemplates(lang) [key]
	if tr == nil {
		tr = GetTrtemplates("default") [key]
	}
	if tr == nil {
		tr, _ = template.New(key).Parse(key)
	}
	return tr
}

func ExpandTemplate(key, lang string, values map[string]interface{}) string {
	var b bytes.Buffer
	tmpl := GetTrtemplate(key, lang)
	err := tmpl.Execute(&b, &values)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func I18n(key, lang string, values ...map[string]interface{})string {
	if len(values) == 0 {
		return strings.TrimSpace(GetTranslation(key, lang))
	} else {
		return strings.TrimSpace(ExpandTemplate(key, lang, values[0]))
	}
}

// regular expressions for filteing HTTP language headers
var cleanSub = regexp.MustCompile("(-[A-Z][A-Z])|(;q=0.\\d)")
var splitter = regexp.MustCompile("(q=0.\\d),")
var quality =	regexp.MustCompile(".*q=0.([0-9])")

// extract languages with q >= 0.5 from HTTP language header
func PreferedLanguages(r *http.Request) []string {
	var langs [4]string
	i := 0
	lnghdr := r.Header["Accept-Language"]
	if len(lnghdr) != 0 && lnghdr[0] != "" {
		slnghdr := strings.Split(splitter.ReplaceAllString(lnghdr[0], "$1|"), "|")
		found := make(map[string]bool)
		outer:        for _, v := range slnghdr {
			qlty, ok := strconv.Atoi(quality.ReplaceAllString(v, "$1"))
			if ok == nil && qlty < 5 {
				break
			}
			lgs := cleanSub.ReplaceAllString(v, "")
			for _, lg := range strings.Split(lgs, ",") {
				if i == 3 {
					break outer
				}
				if ! found[lg] {
					langs[i] = lg
					found[lg] = true
					i++
				}
			}
		}
	}
	if i == 0 {
		langs[0] = DEFAULT_LANGUAGE
		i = 1
	}
	return  langs[0:i]
}
