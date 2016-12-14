package controller

import (
	"reflect"
	"errors"
	"strings"
	"strconv"
)

// type to record value changes for feb forms
type MergeInfo struct {
	Mine     interface{}
	Other    interface{}
	Conflict bool
}

// function MergeDiff runs a three way diff between the old and new versions of my struct and a
// different changed version of the same struct, similar to git merge. If attributes were not changed
// from my old and new versions and automerge is false, differences between me new version and the
// other version are recorded and flagged with conflict = false. If automerge is true, these differences
// are not recorded. Differences between my new changed version and a changed other version are always
// recorded and flagged with conflict = true.
// mineOld, mineNew, other must be pointers to same types
func MergeDiff(mineOld, mineNew, other interface{}, infomerge bool, tags ...string) (diffs map[string]MergeInfo, err error) {
	useTags := len(tags) == 1
	// get value objects of struct variables referenced by interfaces
	vMineOld := reflect.ValueOf(mineOld).Elem()
	vMineNew := reflect.ValueOf(mineNew).Elem()
	vOther := reflect.ValueOf(other).Elem()
	diffs = make(map[string]MergeInfo)
	// make sure you don't compare apples to pears
	if vMineOld.Type() != vOther.Type() || vMineOld.Type() != vMineNew.Type() {
		err = errors.New("different types")
		return
	}
	// loop over all fields
	for i := 0; i < vMineOld.NumField(); i++ {
		fieldInfo := vMineOld.Type().Field(i)
		fieldName := fieldInfo.Name
		fieldTags := fieldInfo.Tag
		fieldMineOld := vMineOld.Field(i).Interface()
		fieldMineNew := vMineNew.Field(i).Interface()
		fieldOther := vOther.Field(i).Interface()
		// ignore anonymous fields
		if !fieldInfo.Anonymous {
			var diffkey string
			if useTags {
				diffkey = fieldTags.Get(tags[0])
				if diffkey == "" {
					continue
				}
			} else {
				diffkey = fieldName
			}
			// other has changed and is different from self -> update self
			if fieldMineNew != fieldOther && fieldOther != fieldMineOld {
				switch fieldInfo.Type.Kind() {
				case reflect.Array:
					fallthrough
				case reflect.Slice:
					mergeArray(vMineOld.Field(i), vMineNew.Field(i), vOther.Field(i),
						diffkey, diffs, infomerge)
				default:
					// if self changed, too -> inform about conflict to be resolved
					// or infomerge is set -> inform about update
					if infomerge || fieldMineNew != fieldMineOld {
						diffs[diffkey] = MergeInfo{fieldMineNew, fieldOther, fieldMineNew != fieldMineOld}
					}
					vMineNew.Field(i).Set(vOther.Field(i))
				}
			}
		}
	}
	return
}

// function mergeArray extends merge behaviour to arrays as parts of structs.
func mergeArray(mold, mnew, other reflect.Value, key string, diffs map[string]MergeInfo, automerge bool) {
	lmo := mold.Len()
	lmn := mnew.Len()
	lot := other.Len()
	idx := 0
	if strings.Contains(key, "#") {
		keyparts := strings.Split(key, "#")
		key = keyparts[0]
		if len(keyparts) > 0 {
			start, err := strconv.Atoi(keyparts[1])
			if err == nil {
				idx = start
			}
		}
	}
	lng := min(lmo, lmn, lot)
	for i := 0; i < lng; i++ {
		vmo := mold.Index(i)
		vmn := mnew.Index(i)
		vot := other.Index(i)
		if vmn.Interface() != vot.Interface() && vmo.Interface() != vot.Interface() {
			// other has changed and is different from self
			if automerge || vmn.Interface() != vmo.Interface() {
				// if self changed, too or automerge is set
				dkey := key + strconv.Itoa(idx + i)
				// record change in map
				diffs[dkey] = MergeInfo{vmn.Interface(), vot.Interface(), vmn.Interface() != vmo.Interface()}
			}
			// update self to changed value
			vmn.Set(vot)
		}
	}

}

func min(i0 int, i ...int) int {
	if len(i) != 0 {
		for _, in := range i {
			if in < i0 {
				i0 = in
			}
		}
	}
	return i0

}
