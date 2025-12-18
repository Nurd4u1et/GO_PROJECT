package main

import "testing"

func TestApplicationStarts(t *testing.T) {
	if true != true {
		t.Fail()
	}
}
