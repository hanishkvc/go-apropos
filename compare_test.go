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
	}{{"st", "testme"}, {"dial", "dialup"}}

	for i := range testData {
		test := testData[i]
		mtp, err := matchtoken_prepare(test.pattern)
		if err != nil {
			t.Errorf("ERRR:Prepare: Mode:%v Pattern:%v Check:%v, Err:%v\n", giMatchMode, test.pattern, test.check, err)
			continue
		}
		ok := mtp.Matchok(test.check)
		if !ok {
			t.Errorf("ERRR:MatchOk: Mode:%v Pattern:%v Check:%v\n", giMatchMode, test.pattern, test.check)
		} else {
			t.Logf("FINE:MatchOk: Mode:%v Pattern:%v Check:%v\n", giMatchMode, test.pattern, test.check)
		}
	}
}
