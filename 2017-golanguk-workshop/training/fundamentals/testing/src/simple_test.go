package main

import "testing"

func TestSimple(t *testing.T) {
	if true {
		t.Error("expected false, got true")
	}
}

/* output
--- FAIL: TestSimple (0.00s)
 simple_test.go:7: expected false, got true
FAIL
*/
