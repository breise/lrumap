package lrumap

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

type exp struct {
	name   string
	n      interface{}
	found  bool
	nItems int
	list   []interface{}
}

type testCase struct {
	max    int
	inputs []exp
}

var intTests = []testCase{
	{5, []exp{
		{"10", 10, false, 1, []interface{}{10}},
		{"20", 20, false, 2, []interface{}{10, 20}},
		{"30", 30, false, 3, []interface{}{10, 20, 30}},
		{"40", 40, false, 4, []interface{}{10, 20, 30, 40}},
		{"50", 50, false, 5, []interface{}{10, 20, 30, 40, 50}},
		{"10", 10, true, 5, []interface{}{20, 30, 40, 50, 10}},
		{"20", 20, true, 5, []interface{}{30, 40, 50, 10, 20}},
		{"30", 30, true, 5, []interface{}{40, 50, 10, 20, 30}},
		{"40", 40, true, 5, []interface{}{50, 10, 20, 30, 40}},
		{"50", 50, true, 5, []interface{}{10, 20, 30, 40, 50}},
	},
	},
	{5, []exp{
		{"10", 10, false, 1, []interface{}{10}},
		{"20", 20, false, 2, []interface{}{10, 20}},
		{"30", 30, false, 3, []interface{}{10, 20, 30}},
		{"40", 40, false, 4, []interface{}{10, 20, 30, 40}},
		{"50", 50, false, 5, []interface{}{10, 20, 30, 40, 50}},
		{"60", 60, false, 5, []interface{}{20, 30, 40, 50, 60}},
		{"20", 20, true, 5, []interface{}{30, 40, 50, 60, 20}},
		{"30", 30, true, 5, []interface{}{40, 50, 60, 20, 30}},
		{"40", 40, true, 5, []interface{}{50, 60, 20, 30, 40}},
		{"50", 50, true, 5, []interface{}{60, 20, 30, 40, 50}},
		{"60", 60, true, 5, []interface{}{20, 30, 40, 50, 60}},
		{"10", 10, false, 5, []interface{}{30, 40, 50, 60, 10}},
		{"20", 20, false, 5, []interface{}{40, 50, 60, 10, 20}},
	},
	},
}

func TestBasic(t *testing.T) {
	for i, tc := range intTests {
		desc := fmt.Sprintf("Test Case %d:", i)
		t.Run(desc, func(t *testing.T) {
			singleCaseTest(t, tc, desc, sumTo)
		})
	}
}

var structTests = []string{
	"testdata/shakespeare_plays/alls-well-that-ends-well_TXT_FolgerShakespeare.txt",
	"testdata/shakespeare_plays/coriolanus_TXT_FolgerShakespeare.txt",
	"testdata/shakespeare_plays/hamlet_TXT_FolgerShakespeare.txt",
	"testdata/shakespeare_plays/julius-caesar_TXT_FolgerShakespeare.txt",
}

type play struct {
	title, author, editors1, editors2, source, url, createdOn, text string
}

func TestComplex(t *testing.T) {

	plays := make([]play, len(structTests))

	for i, filePath := range structTests {
		txt, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("cannot open file %s for reading: %s", filePath, err)
		}
		info := strings.SplitN(string(txt), "\n", 8)
		plays[i] = play{
			title:     info[0], // All's Well That Ends Well
			author:    info[1], // by William Shakespeare
			editors1:  info[2], // Edited by Barbara A. Mowat and Paul Werstine
			editors2:  info[3], //   with Michael Poston and Rebecca Niles
			source:    info[4], // Folger Shakespeare Library
			url:       info[5], // https://shakespeare.folger.edu/shakespeares-works/alls-well-that-ends-well/
			createdOn: info[6], // Created on Mar 14, 2018, from FDT version 0.9.2.2
			text:      info[7],
		}
	}
	var exps0 []exp
	for i, play := range plays {
		nm := fmt.Sprintf("%s_0_%d", strings.TrimSpace(play.title), i)
		exps0 = append(exps0, exp{name: nm, n: play, found: false, nItems: i + 1})
	}
	for i := range plays {
		idx := len(plays) - 1 - i
		play := plays[idx]
		nm := fmt.Sprintf("%s_1_%d", strings.TrimSpace(play.title), i)
		exps0 = append(exps0, exp{name: nm, n: play, found: true, nItems: len(plays)})
	}
	tc0 := testCase{max: len(plays), inputs: exps0}

	var exps1 []exp
	maxItems := len(plays) - 1
	for i, play := range plays {
		nm := fmt.Sprintf("%s_0_%d", strings.TrimSpace(play.title), i)
		exps1 = append(exps1, exp{name: nm, n: play, found: false, nItems: minOf(i+1, maxItems)})
	}
	for i := range plays {
		idx := len(plays) - 1 - i
		play := plays[idx]
		nm := fmt.Sprintf("%s_1_%d", strings.TrimSpace(play.title), i)
		exps1 = append(exps1, exp{name: nm, n: play, found: idx > 0, nItems: maxItems})
	}
	tc1 := testCase{max: maxItems, inputs: exps1}

	theseTests := []testCase{
		tc0,
		tc1,
	}
	for i, tc := range theseTests {
		desc := fmt.Sprintf("Test Case %d:", i)
		t.Run(desc, func(t *testing.T) {
			singleCaseTest(t, tc, desc, countWords)
		})
	}

}

func singleCaseTest(t *testing.T, tc testCase, desc string, fn func(interface{}) interface{}) {
	lm := New().MaxItems(tc.max)
	for j, v := range tc.inputs {
		msgFmt := "%s: input #%d: %s: %s: Exp: %v; Got: %v"
		got, found := lm.Get(v.n)
		if found != v.found {
			t.Errorf(msgFmt, desc, j, v.name, "found", v.found, found)
		}
		if !found {
			lm.Put(v.n, fn(v.n))
		}
		if lm.NItems() != v.nItems {
			t.Errorf(msgFmt, desc, j, v.name, "nItems", v.nItems, lm.NItems())
		}
		if found && got != fn(v.n) {
			t.Errorf(msgFmt, desc, j, v.name, "got", fn(v.n), got)
		}
		// if found && got == fn(v.n) {
		// 	fmt.Printf(msgFmt+"\n", desc, j, v, "got", fn(v.n), got)
		// }
		if len(v.list) > 0 {
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
				t.Errorf(msgFmt, desc, j, v.name, "contents", expSl, gotSl)
			}
		}
	}
}
func sumTo(x interface{}) interface{} {
	n, ok := x.(int)
	if !ok {
		panic(fmt.Sprintf("cannot cast %v as int", x))
	}
	var rv int
	for i := 0; i < n; i++ {
		rv += i
	}
	return rv
}

func countWords(x interface{}) interface{} {
	n, ok := x.(play)
	if !ok {
		panic(fmt.Sprintf("cannot cast %v as play", x))
	}
	rv := strings.Split(n.text, " ")
	return len(rv)
}

func minOf(x, y int) int {
	if x < y {
		return x
	}
	return y
}
