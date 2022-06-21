// Test main db related logics
// HanishKVC, 2022

package main

import (
	"testing"
)

func TestFind(t *testing.T) {
	gbSortedResult = true
	load_dbs()
	for _, mm := range []MatchMode{MatchMode_Contains, MatchMode_RegExp} {
		for _, cs := range []bool{true, false} {
			t.Logf("MatchMode:%v CaseSensitive:%v\n", matchmode_tostr(mm), cs)
			db_find(gDB, "fmt", FINDCMT_DUMMY, FINDPKG_DEFAULT, mm, cs)
		}
	}
}

func BenchmarkFindContains(b *testing.B) {
	gbSortedResult = true
	gbCaseSensitive = false
	for i := 0; i < b.N; i++ {
		load_dbs()
		db_find(gDB, "numcpu", FINDCMT_DUMMY, FINDPKG_DEFAULT, MatchMode_Contains, gbCaseSensitive)
	}
}

func BenchmarkFindRegexp(b *testing.B) {
	gbSortedResult = true
	gbCaseSensitive = false
	for i := 0; i < b.N; i++ {
		load_dbs()
		db_find(gDB, "numcpu", FINDCMT_DUMMY, FINDPKG_DEFAULT, MatchMode_RegExp, gbCaseSensitive)
	}
}
