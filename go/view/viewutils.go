package view

import (
	"errors"
	"reflect"
	"fmt"
	"strconv"
	"html/template"
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
		return nil, errors.New("invalid Dict call")
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


// function DotReference helps to pass a field reference of Dot {{.}} or some other
// variable as a parameter into a nested template
// params: dot the template Dot or some other variable of type struct, map, array or slice
//	   at  the field name, map key or slice index you want to address
// returns the addressed value as string or an error if not applicable
func DotReference(dot interface{}, at string) (result string, err error) {
	var value interface{}
	reflected := reflect.ValueOf(dot)
	switch reflected.Kind() {
	case reflect.Struct:
		value = reflected.FieldByName(at).String()
		err = nil
	case reflect.Map:
		value = reflected.MapIndex(reflect.ValueOf(at))
		err = nil
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		idx, e := strconv.Atoi(at)
		if e == nil && idx < reflected.Len() {
			value = reflected.Index(idx)
			err = nil
		} else if e != nil {
			err = e
			return
		} else {
			err = errors.New("index out of range")
			return
		}
	default:
		err = errors.New(fmt.Sprintf("%s in not supported", reflected.Kind().String()))
		return
	}
	result = fmt.Sprintf("%v", value)
	return
}
