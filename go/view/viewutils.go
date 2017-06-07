package view

import (
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

// Function Dict creates a map from its inputs for use in go templates
// thanks to stackoverflow user tux21b
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values) % 2 != 0 {
		return nil, errors.New("invalid Dict call")
	}
	dict := make(map[string]interface{}, len(values) / 2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i + 1]
	}
	return dict, nil
}

// function SaveAtt flags the string representation of an
// interface parameter as a safe html attribute
func SafeAtt(s interface{}) template.HTMLAttr {
	if s != nil {
		return template.HTMLAttr(fmt.Sprint(s))
	}
	return template.HTMLAttr("")
}

// function Concat cocatenates the string representations of its parameters
func Concat(s ...interface{}) string {
	r := ""
	for _, e := range s {
		r += fmt.Sprint(e)
	}
	return r
}

// function AddDict adds more key/value pairs to a dictionary. This helps to break long parameter
// lists to templates on several lines to make them better readable
func AddDict(dict map[string]interface{}, values ...interface{}) (map[string]interface{}, error) {
	if len(values) % 2 != 0 {
		return nil, errors.New("invalid AddDict call")
	}
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i + 1]
	}
	return dict, nil
}

// function MergeDict merges dictionaries for use in templates into
// a new dictionary
func MergeDict(dict map[string]interface{}, more ...map[string]interface{}) (map[string]interface{}, error) {
	if len(more) == 0 {
		return nil, errors.New("invalid MergeDict call")
	}
	l:= len(dict)
	for _, m := range more {
		l += len(m)
	}
	ret := make(map[string] interface{}, l)
	for k, v := range dict {
		ret[k] = v
	}
	for _, m := range more {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret, nil
}

// function IsKind tests an object obj for a given reflect.Kind
func IsKind(obj interface{}, kind string) bool {
	return strings.ToUpper(reflect.TypeOf(obj).Kind().String()) == strings.ToUpper(kind)
}

func IsMod(val, mod int) bool {
	return val > 0 && val % mod == 0
}