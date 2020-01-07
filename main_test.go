package main

import (
	"io/ioutil"
	"os"
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

func TestSimpleQuotes(t *testing.T) {
	rawf := `
-- define: test1
aa bb "--cc" "-d--dd"
ee "--" ff
-- end
	`
	mk := make([]string, 0)
	mv := make([]string, 0)
	mtags := make([][]string, 0)
	f, err := ioutil.TempFile("", "simplequotes")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	fnm := f.Name()
	_, _ = f.Write([]byte(rawf))
	f.Close()
	defer os.Remove(fnm)
	extractall([]string{fnm}, &mk, &mv, &mtags)
	if len(mk) != 1 || len(mv) != 1 {
		t.Fatal("invalid len")
	}
	if mv[0] != "aa bb \"--cc\" \"-d--dd\"\nee \"--\" ff\n" {
		t.Fatal("invalid chars", mv[0])
	}
}
