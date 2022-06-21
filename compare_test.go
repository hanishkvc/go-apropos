// Test some compare helpers
// HanishKVC, 2022

package main

import (
	"regexp"
	"testing"
)

func TestMatcher(t *testing.T) {
	var testData = []struct {
		pattern string
		check   string
		expect  bool
	}{{"st", "testme", true}, {"dial", "dialup", true}, {"dial", "testme", false}}

	for _, caseSensitive := range []bool{true, false} {
		for _, m := range []MatchMode{MatchMode_Contains, MatchMode_RegExp} {
			sMatchMode := matchmode_tostr(m)
			for i := range testData {
				test := testData[i]
				matcher := New_Matcher(m, test.pattern, caseSensitive)
				ok := matcher.Matchok(test.check)
				if ok != test.expect {
					t.Errorf("ERRR:MatchOk: Mode:%v[%v] Pattern:%v[%v] Check:%v Expected:%v Got:%v\n", sMatchMode, matcher.MType(), test.pattern, matcher.Pattern(), test.check, test.expect, ok)
				} else {
					t.Logf("FINE:MatchOk: Mode:%v[%v] Pattern:%v[%v] Check:%v Expected:%v Got:%v\n", sMatchMode, matcher.MType(), test.pattern, matcher.Pattern(), test.check, test.expect, ok)
				}
			}
		}
	}
}

const BMRE_PATTERN = "st"
const BMRE_CHECK = "testme"

func BenchmarkRECompile(b *testing.B) {
	re := regexp.MustCompile(BMRE_PATTERN)
	for i := 0; i < b.N; i++ {
		re.MatchString(BMRE_CHECK)
	}
}

func BenchmarkRENoCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		regexp.MatchString(BMRE_PATTERN, BMRE_CHECK)
	}
}
