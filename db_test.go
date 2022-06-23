// Test main db related logics
// HanishKVC, 2022

package main

import (
	"path/filepath"
	"testing"
)

func TestFind(t *testing.T) {
	sortedResult := true
	sFind := "Fmt"
	load_dbs()
	for _, mm := range []MatchMode{MatchMode_Contains, MatchMode_RegExp} {
		for _, cs := range []bool{true, false} {
			t.Logf("Find: %v MatchMode:%v CaseSensitive:%v\n", sFind, matchmode_tostr(mm), cs)
			db_find(gDB, sFind, FINDCMT_DUMMY, FINDPKG_DEFAULT, mm, cs, sortedResult)
		}
	}
}

func BenchmarkFindContains(b *testing.B) {
	sortedResult := true
	caseSensitive := false
	for i := 0; i < b.N; i++ {
		load_dbs()
		db_find(gDB, "numcpu", FINDCMT_DUMMY, FINDPKG_DEFAULT, MatchMode_Contains, caseSensitive, sortedResult)
	}
}

func BenchmarkFindRegexp(b *testing.B) {
	sortedResult := true
	caseSensitive := false
	for i := 0; i < b.N; i++ {
		load_dbs()
		db_find(gDB, "numcpu", FINDCMT_DUMMY, FINDPKG_DEFAULT, MatchMode_RegExp, caseSensitive, sortedResult)
	}
}

func TestPkgBasePath(t *testing.T) {
	srcBPath := filepath.Join("root", "srcbpath")
	aTests := []struct {
		pkgName string
		sFile   string
	}{
		{"pkg", filepath.Join(srcBPath, "pkg", "pkg.go")},
		{"pkg", filepath.Join(srcBPath, "pkg", "other.go")},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "pkg.go")},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "other.go")},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "pkg-other.go")},
	}
	for _, c := range aTests {
		t.Logf("%v: %v:%v [%v]", srcBPath, c.pkgName, c.sFile, pkg_basepath(c.pkgName, c.sFile, srcBPath, true))
	}
}
