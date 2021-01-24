package lrumap

import (
	"fmt"
	"testing"
)

type exp struct {
	n      int
	found  bool
	nItems int
	list   []int
}

var tests = []struct {
	max    int
	inputs []exp
}{
	{5, []exp{
		{10, false, 1, []int{10}},
		{20, false, 2, []int{10, 20}},
		{30, false, 3, []int{10, 20, 30}},
		{40, false, 4, []int{10, 20, 30, 40}},
		{50, false, 5, []int{10, 20, 30, 40, 50}},
		{10, true, 5, []int{20, 30, 40, 50, 10}},
		{20, true, 5, []int{30, 40, 50, 10, 20}},
		{30, true, 5, []int{40, 50, 10, 20, 30}},
		{40, true, 5, []int{50, 10, 20, 30, 40}},
		{50, true, 5, []int{10, 20, 30, 40, 50}},
	},
	},
	{5, []exp{
		{10, false, 1, []int{10}},
		{20, false, 2, []int{10, 20}},
		{30, false, 3, []int{10, 20, 30}},
		{40, false, 4, []int{10, 20, 30, 40}},
		{50, false, 5, []int{10, 20, 30, 40, 50}},
		{60, false, 5, []int{20, 30, 40, 50, 60}},
		{20, true, 5, []int{30, 40, 50, 60, 20}},
		{30, true, 5, []int{40, 50, 60, 20, 30}},
		{40, true, 5, []int{50, 60, 20, 30, 40}},
		{50, true, 5, []int{60, 20, 30, 40, 50}},
		{60, true, 5, []int{20, 30, 40, 50, 60}},
		{10, false, 5, []int{30, 40, 50, 60, 10}},
		{20, false, 5, []int{40, 50, 60, 10, 20}},
	},
	},
}

func TestBasic(t *testing.T) {
	for i, tc := range tests {
		desc := fmt.Sprintf("Test Case %d: %v", i, tc)
		t.Run(desc, func(t *testing.T) {
			lm := New().MaxItems(tc.max)
			for j, v := range tc.inputs {
				msgFmt := "%s: input #%d: %v: %s: Exp: %v; Got: %v"
				got, found := lm.Get(v.n)
				if found != v.found {
					t.Errorf(msgFmt, desc, j, v, "found", v.found, found)
				}
				if !found {
					lm.Put(v.n, sumTo(v.n))
				}
				if lm.NItems() != v.nItems {
					t.Errorf(msgFmt, desc, j, v, "nItems", v.nItems, lm.NItems())
				}
				if found && got != sumTo(v.n) {
					t.Errorf(msgFmt, desc, j, v, "got", sumTo(v.n), got)
				}
				kvSl := lm.lruList.ToSlice()
				kSl := make([]interface{}, len(kvSl))
				for k, kv := range kvSl {
					kvp, ok := kv.(kvPair)
					if !ok {
						t.Fatalf("%s: input #%d: %v: ToSlice(): element %d cannot be cast to kvPair (%v)", desc, j, v, k, kv)
					}
					kSl[k] = kvp.k
				}
				gotSl := fmt.Sprintf("%v", kSl)
				expSl := fmt.Sprintf("%v", v.list)
				if gotSl != expSl {
					t.Errorf(msgFmt, desc, j, v, "contents", expSl, gotSl)
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
