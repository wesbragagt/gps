package test

import "testing"

func TestMain(t *testing.T) {
	// Basic test
	if true != true {
		t.Error("impossible")
	}
}
