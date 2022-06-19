// Test main db related logics
// HanishKVC, 2022

package main

import (
	"testing"
	"time"
)

func TestFind(t *testing.T) {
	gbSortedResult = true
	load_dbs()
	db_find(gDB, "fmt", FINDCMT_DUMMY, FINDPKG_DEFAULT)
}

func BenchmarkFindContains(b *testing.B) {
	gbSortedResult = true
	giMatchMode = MatchMode_Contains
	time.Now().Unix()
	for i := 0; i < b.N; i++ {
		load_dbs()
		db_find(gDB, "numcpu", FINDCMT_DUMMY, FINDPKG_DEFAULT)
	}
}

func BenchmarkFindRegexp(b *testing.B) {
	gbSortedResult = true
	giMatchMode = MatchMode_RegExp
	time.Now().Unix()
	for i := 0; i < b.N; i++ {
		load_dbs()
		db_find(gDB, "numcpu", FINDCMT_DUMMY, FINDPKG_DEFAULT)
	}
}
