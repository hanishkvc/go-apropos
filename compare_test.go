// Test some compare helpers
// HanishKVC, 2022

package main

import (
	"testing"
)

func TestMatch(t *testing.T) {
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
			mtp, err := matchtoken_prepare(test.pattern)
			if err != nil {
				t.Errorf("ERRR:Prepare: Mode:%v Pattern:%v Check:%v, Err:%v\n", sMatchMode, test.pattern, test.check, err)
				continue
			}
			ok := mtp.Matchok(test.check)
			if ok != test.expect {
				t.Errorf("ERRR:MatchOk: Mode:%v Pattern:%v Check:%v Expected:%v Got:%v\n", sMatchMode, test.pattern, test.check, test.expect, ok)
			} else {
				t.Logf("FINE:MatchOk: Mode:%v Pattern:%v Check:%v Expected:%v Got:%v\n", sMatchMode, test.pattern, test.check, test.expect, ok)
			}
		}
	}
}
