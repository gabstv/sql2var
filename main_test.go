package main

import (
	"testing"
)

func TestSlice(t *testing.T) {
	v := "[]string{\"a\", \"b\"}"
	raw := `a


	;
b `
	v2 := getslicecode(raw, ";")
	if v != v2 {
		t.Fatal(v, v2)
	}
}
