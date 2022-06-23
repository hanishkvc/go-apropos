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
		expect  string
	}{
		{"pkg", filepath.Join(srcBPath, "pkg", "pkg.go"), "pkg"},
		{"pkg", filepath.Join(srcBPath, "pkg", "other.go"), "pkg"},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "pkg.go"), "xbasepath/pkg"},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "other.go"), "xbasepath/pkg"},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "pkg", "pkg.go"), "xbasepath/pkg"},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "pkg", "other.go"), "xbasepath/pkg"},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "pkg-other.go"), "xbasepath/pkg"},
		{"pkg", filepath.Join(srcBPath, "xbasepath", "other-pkg.go"), "xbasepath/pkg"},
	}
	for _, c := range aTests {
		sPreGot := pkg_basepath_preslash(c.pkgName, c.sFile, srcBPath, true)
		if sPreGot == c.expect {
			t.Logf("INFO:PRE :%v: %v:%v [%v = %v]", srcBPath, c.pkgName, c.sFile, sPreGot, c.expect)
		} else {
			t.Errorf("ERRR:PRE :%v: %v:%v [%v != %v]", srcBPath, c.pkgName, c.sFile, sPreGot, c.expect)
		}
	}
	for _, c := range aTests {
		sPostGot := pkg_basepath_postslash(c.pkgName, c.sFile, srcBPath, true)
		if sPostGot == c.expect {
			t.Logf("INFO:POST:%v: %v:%v [%v = %v]", srcBPath, c.pkgName, c.sFile, sPostGot, c.expect)
		} else {
			t.Errorf("ERRR:POST:%v: %v:%v [%v != %v]", srcBPath, c.pkgName, c.sFile, sPostGot, c.expect)
		}
	}
}

func BenchmarkPkgBasePathPre(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pkg_basepath_preslash("pkg", "/root/srcpath/xbasepath/pkg/pkg.go", "/root/srcpath", true)
		pkg_basepath_preslash("pkg", "/root/srcpath/xbasepath/abc/other.go", "/root/srcpath", true)
	}
}

func BenchmarkPkgBasePathPost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pkg_basepath_postslash("pkg", "/root/srcpath/xbasepath/pkg/pkg.go", "/root/srcpath", true)
		pkg_basepath_postslash("pkg", "/root/srcpath/xbasepath/abc/other.go", "/root/srcpath", true)
	}
}
