package controller

import (
	"reflect"
	"errors"
	"strings"
	"strconv"
	"fmt"
)

type MergeDiffType int

const (
	NONE MergeDiffType = iota       // no changes at all
	MINE                            // only my value changed
	THEIRS                          // only their value changed
	BOTH                            // both values changed and are different
	SAME                            // both values changed but are equal
)

// type to record value changes for web forms
type MergeInfo struct {
	Mine     interface{}
	Other    interface{}
	Conflict MergeDiffType
}

// function MergeDiff runs a three way diff between the old and new versions of my struct and a
// different changed version of the same struct, similar to git merge. Differences are flagged in
// return parameter diffs using type MergeDiffType.If automerge is true, mineNew fields get
// updated to updated fields of otherNew and only changed fields are flagged in diffs.
// mineOld, mineNew, otherNew must be pointers to same struct types.
// If tag is not empty, only fields are compared that are tagged with the given tag and the tag value
// will be used as key in the diffs map.
func MergeDiff(mineOld, mineNew, otherNew interface{}, automerge bool, tag ...string) (diffs map[string]MergeInfo, err error) {
	useTags := len(tag) == 1
	// get value objects of struct variables referenced by interfaces
	vMineOld := reflect.ValueOf(mineOld).Elem()
	vMineNew := reflect.ValueOf(mineNew).Elem()
	vOther := reflect.ValueOf(otherNew).Elem()
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
				diffkey = fieldTags.Get(tag[0])
				if diffkey == "" {
					continue
				}
			} else {
				diffkey = fieldName
			}
			// simple field or array/slice?
			switch fieldInfo.Type.Kind() {
			case reflect.Array:
				fallthrough
			case reflect.Slice:
				mergeArray(vMineOld.Field(i), vMineNew.Field(i), vOther.Field(i),
					diffkey, diffs, automerge)
			default:
				cmine := fieldMineNew == fieldMineOld   // has my field changed
				ctheirs := fieldOther == fieldMineOld   // has their field changed
				cdiff := fieldMineNew == fieldOther     // are new fields same
				var mergeDiff MergeDiffType
				if cdiff && cmine {
					// no change
					mergeDiff = NONE
				} else if !cdiff && cmine {
					// they changed
					mergeDiff = THEIRS
				} else if !cdiff && !cmine && ctheirs {
					// mine changed
					mergeDiff = MINE
				} else if !cdiff && !cmine && ! ctheirs {
					// both changed differently
					mergeDiff = BOTH
				} else {
					// both changed but are same
					mergeDiff = SAME
				}
				switch mergeDiff {
				case THEIRS:
					fallthrough
				case BOTH:
					//vMineNew.Field(i).Set(vOther.Field(i))
					diffs[diffkey] = MergeInfo{fieldMineNew, fieldOther, mergeDiff}
				case MINE:
					fallthrough
				case SAME:
					if automerge {
						diffs[diffkey] = MergeInfo{fieldMineNew, fieldOther, mergeDiff}
					}
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
		vmo := mold.Index(i).Interface()
		vmn := mnew.Index(i).Interface()
		vot := other.Index(i).Interface()
		cmine := vmn == vmo     // has my field changed
		ctheirs := vot == vmo   // has their field changed
		cdiff := vmn == vot     // are new fields same
		var mergeDiff MergeDiffType
		if cdiff && cmine {
			// no change
			mergeDiff = NONE
		} else if !cdiff && cmine {
			// they changed
			mergeDiff = THEIRS
		} else if !cdiff && !cmine && ctheirs {
			// mine changed
			mergeDiff = MINE
		} else if !cdiff && !cmine && ! ctheirs {
			// both changed differently
			mergeDiff = BOTH
		} else {
			// both changed but are same
			mergeDiff = SAME
		}
		dkey := key + strconv.Itoa(idx + i)
		switch mergeDiff {
		case THEIRS:
			fallthrough
		case BOTH:
			mnew.Index(i).Set(other.Index(i))
			diffs[dkey] = MergeInfo{vmn, vot, mergeDiff}
		case MINE:
			fallthrough
		case SAME:
			if automerge {
				diffs[dkey] = MergeInfo{vmn, vot, mergeDiff}
			}
		}
	}

}

// function MergeScale scales integral values by factor for all
// map keys that contain string key
func MergeScaleResults(diffs map[string]MergeInfo, key string) {
	for k, v := range diffs {
		if strings.Contains(k, key) && reflect.ValueOf(v.Mine).Kind() == reflect.Int {
			diffs[k] = MergeInfo{
				fmt.Sprintf("%.1f", float32(v.Mine.(int)) / 10.),
				fmt.Sprintf("%.1f", float32(v.Other.(int)) / 10.),
				v.Conflict,
			}
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
