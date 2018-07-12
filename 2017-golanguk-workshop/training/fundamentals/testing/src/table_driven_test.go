package main

import "testing"

func TestTableDriven(t *testing.T) {
	tt := []struct {
		A        int
		B        int
		Expected int
	}{
		{A: 1, B: 1, Expected: 2},
		{A: 2, B: 2, Expected: 4},
		{A: 3, B: 3, Expected: 5},
		{A: 4, B: 4, Expected: 6},
	}

	for _, x := range tt {
		got := x.A + x.B
		if got != x.Expected {
			t.Errorf("expected %d, got %d", x.Expected, got)
		}
	}
}

/* output
--- FAIL: TestTableDriven (0.00s)
 table_driven_test.go:20: expected 5, got 6
 table_driven_test.go:20: expected 6, got 8
*/
