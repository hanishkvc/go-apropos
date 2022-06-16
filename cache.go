// Helper routines wrt Cache
// HanishKVC, 2022
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var gsGoVersion = ""
var gCacheBase = "~/.cache"

const gDBAllCacheFile = "goapropos.dball"
const gsFNCacheVersion = "goapropos.ver"

func cache_filename(cacheFile string) (string, error) {
	sCacheBase, err := adjust_path(gCacheBase)
	if err != nil {
		return "", err
	}
	sDBCacheFile := sCacheBase + string(os.PathSeparator) + cacheFile
	return sDBCacheFile, nil
}

func cache_force_fresh() {
	gbUseCache = false
	gbCreateCache = true
	if gFindPkg != FINDPKG_DEFAULT {
		gFindPkg = FINDPKG_DEFAULT
		fmt.Printf("%v:WARN:Cache: Ignoring FindPkg\n", PRG_TAG)
	}
	fmt.Printf("%v:INFO:Cache: Will create/update cache...\n", PRG_TAG)
}

func cache_writever() {
	fName, err := cache_filename(gsFNCacheVersion)
	if err != nil {
		fmt.Printf("%v:ERRR:Cache: WriteVer:Filename:%v\n", PRG_TAG, err)
		return
	}
	err = os.WriteFile(fName, []byte(gsGoVersion), 0600)
	if err != nil {
		fmt.Printf("%v:ERRR:Cache: WriteVer:WriteFile:%v\n", PRG_TAG, err)
	}
	if giDEBUG > -1 {
		fmt.Printf("%v:Info:Cache: WriteVer:Done\n", PRG_TAG)
	}
}

func cache_ok_or_fresh() {
	fName, err := cache_filename(gsFNCacheVersion)
	if err != nil {
		fmt.Printf("%v:ERRR:Cache: Filename:%v\n", PRG_TAG, err)
		cache_force_fresh()
		return
	}
	bsCacheVer, err := os.ReadFile(fName)
	if err != nil {
		fmt.Printf("%v:ERRR:Cache: VersionFile:%v\n", PRG_TAG, err)
		cache_force_fresh()
		return
	}
	sCacheVer := strings.TrimSpace(string(bsCacheVer))
	if sCacheVer != gsGoVersion {
		fmt.Printf("%v:WARN:Cache: Version mismatch [%v != %v]\n", PRG_TAG, sCacheVer, gsGoVersion)
		cache_force_fresh()
	}
}

func cache_maya() {
	gsGoVersion = runtime.Version()
	cache_ok_or_fresh()
}
