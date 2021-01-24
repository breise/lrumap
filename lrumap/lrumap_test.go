package lrumap

import (
	"fmt"
	"testing"
)

type exp struct {
	n      int
	found  bool
	nItems int
}

var tests = []struct {
	max    int
	inputs []exp
}{
	{5, []exp{
		{10, false, 1},
		{20, false, 2},
		{30, false, 3},
		{40, false, 4},
		{50, false, 5},
		{10, true, 5},
		{20, true, 5},
		{30, true, 5},
		{40, true, 5},
		{50, true, 5},
	},
	},
	{5, []exp{
		{10, false, 1},
		{20, false, 2},
		{30, false, 3},
		{40, false, 4},
		{50, false, 5},
		{60, false, 5},
		{20, true, 5},
		{30, true, 5},
		{40, true, 5},
		{50, true, 5},
		{60, true, 5},
		{10, false, 5},
		{20, false, 5},
	},
	},
}

func TestBasic(t *testing.T) {
	for i, tc := range tests {
		desc := fmt.Sprintf("Test Case %d: %v", i, tc)
		t.Run(desc, func(t *testing.T) {
			lm := New().MaxItems(tc.max)
			for j, v := range tc.inputs {
				got, found := lm.Get(v.n)
				if found != v.found {
					t.Errorf("%s: input #%d: %v: found: Exp: %v; Got: %v", desc, j, v, v.found, found)
				}
				if !found {
					lm.Put(v.n, sumTo(v.n))
				}
				if lm.NItems() != v.nItems {
					t.Errorf("%s: input #%d: %v: nItems: Exp: %v; Got: %v", desc, j, v, v.nItems, lm.NItems())
				}
				if found && got != sumTo(v.n) {
					t.Errorf("%s: input #%d: %v: got: Exp: %v; Got: %v", desc, j, v, sumTo(v.n), got)
				}
			}
		})
	}
}

func sumTo(n int) int {
	var rv int
	for i := 0; i < n; i++ {
		rv += i
	}
	return rv
}
