package test_util

import (
	"reflect"
	"runtime/debug"
	"testing"
)

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, found, expected interface{}) {
	if found == expected {
		return
	}
	// debug.PrintStack()
	t.Errorf(string(debug.Stack()))
	t.Errorf("Received %v (type %v), expected %v (type %v)", found, reflect.TypeOf(found), expected, reflect.TypeOf(expected))
	t.FailNow()
}

func AssertNotNull(t *testing.T, found interface{}) {
	if found != nil {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received nil")
	t.FailNow()
}
