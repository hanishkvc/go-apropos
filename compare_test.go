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

	for _, m := range []MatchMode{MatchMode_Contains, MatchMode_RegExp} {
		giMatchMode = m
		sMatchMode := matchmode_tostr(giMatchMode)
		for i := range testData {
			test := testData[i]
			mtp, err := matcher_create(test.pattern)
			if err != nil {
				t.Errorf("ERRR:Prepare: Mode:%v Pattern:%v Check:%v, Err:%v\n", sMatchMode, test.pattern, test.check, err)
				continue
			}
			ok := mtp.Matchok(test.check)
			if ok != test.expect {
				t.Errorf("ERRR:MatchOk: Mode:%v[%v] Pattern:%v Check:%v Expected:%v Got:%v\n", sMatchMode, mtp.Utype(), test.pattern, test.check, test.expect, ok)
			} else {
				t.Logf("FINE:MatchOk: Mode:%v[%v] Pattern:%v Check:%v Expected:%v Got:%v\n", sMatchMode, mtp.Utype(), test.pattern, test.check, test.expect, ok)
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
