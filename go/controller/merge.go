package controller

import (
	"reflect"
	"errors"
)

// type to record value changes for feb forms
type Merge struct {
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
func MergeDiff(mineOld, mineNew, other interface{}, automerge bool) (diffs map[string]Merge, err error) {
	vMineOld := reflect.ValueOf(mineOld).Elem()
	vMineNew := reflect.ValueOf(mineNew).Elem()
	vOther := reflect.ValueOf(other).Elem()
	diffs = make(map[string]Merge)
	if vMineOld.Type() != vOther.Type() || vMineOld.Type() != vMineNew.Type() {
		err = errors.New("different types")
		return
	}
	for i := 0; i < vMineOld.NumField(); i++ {
		fieldInfo := vMineOld.Type().Field(i)
		fieldName := fieldInfo.Name
		fieldMineOld := vMineOld.Field(i).Interface()
		fieldMineNew := vMineNew.Field(i).Interface()
		fieldOther := vOther.Field(i).Interface()
		if !fieldInfo.Anonymous {
			if fieldMineNew != fieldMineOld && fieldOther != fieldMineOld && fieldMineNew != fieldOther {
				// both changed
				diffs[fieldName] = Merge{fieldMineNew, fieldOther, true}
				vMineNew.Field(i).Set(vOther.Field(i))
			} else if fieldMineNew == fieldMineOld && fieldMineNew != fieldOther {
				// only other changed
				vMineNew.Field(i).Set(vOther.Field(i))
				if !automerge {
					diffs[fieldName] = Merge{fieldMineNew, fieldOther, false}
				}
			}
		}
	}
	return
}

