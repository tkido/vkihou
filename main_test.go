package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	got := convert("./_testdata/source.txt")

	f, err := os.Open("./_testdata/result.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	want := string(bs)
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
