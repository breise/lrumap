package lrulist

import (
	"fmt"
	"testing"
)

var tests = []struct {
	max, n, expN int
}{
	{10, 10, 10},
	{10, 15, 10},
	{20, 10, 10},
}

func TestLruList(t *testing.T) {
	for i, tc := range tests {
		desc := fmt.Sprintf("Test Case %d: %v", i, tc)
		t.Run(desc, func(t *testing.T) {
			l := New().MaxItems(tc.max)
			for i := 0; i < tc.n; i++ {
				node, dropped := l.Add(i)
				fmt.Printf("on iteration %d, added node: %v; dropped items: %v\n", i, node, dropped)
			}
			fmt.Printf("contents: %+v\n", l.ToSlice())
			if l.NItems() != tc.expN {
				t.Errorf("added %d, expected %d items, but got %d items", tc.n, tc.expN, l.NItems())
			}
		})
	}

}