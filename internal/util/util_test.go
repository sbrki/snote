package util

import (
	"testing"
)

func TestSliceContainsString(t *testing.T) {
	s := make([]string, 3)
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"

	if SliceContainsString(s, "b") != true {
		t.Error("wrong return value")
	}

	if SliceContainsString(s, "d") != false {
		t.Error("wrong return value")
	}
}

func TestSliceRemoveString(t *testing.T) {
	// test 1: removing when one instance of target string is present
	a := make([]string, 4)
	a[0] = "a"
	a[1] = "b"
	a[2] = "b"
	a[3] = "c"

	SliceRemoveString(&a, "a")
	if len(a) != 3 || a[0] != "b" || a[1] != "b" || a[2] != "c" {
		t.Error("wrong return value")
	}

	// test 2: removing when duplicates of target string are present
	a = make([]string, 4)
	a[0] = "a"
	a[1] = "b"
	a[2] = "b"
	a[3] = "c"

	SliceRemoveString(&a, "b")
	if len(a) != 2 || a[0] != "a" || a[1] != "c" {
		t.Error("wrong return value")
	}

	// test 3: removing when target string is not present
	a = make([]string, 4)
	a[0] = "a"
	a[1] = "b"
	a[2] = "b"
	a[3] = "c"

	SliceRemoveString(&a, "d")
	if len(a) != 4 || a[0] != "a" || a[1] != "b" || a[2] != "b" || a[3] != "c" {
		t.Error("wrong return value")
	}
}
