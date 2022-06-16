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

func cache_force_fresh() {
	gbUseCache = false
	gbCreateCache = true
	if gFindPkg != FINDPKG_DEFAULT {
		gFindPkg = FINDPKG_DEFAULT
		fmt.Printf("%v:WARN:Cache: Ignoring FindPkg\n", PRG_TAG)
	}
	fmt.Printf("%v:INFO:Cache: Will create/update cache...\n", PRG_TAG)
}

func cache_ok_or_fresh() {
	bsCacheVer, err := os.ReadFile(gsFNCacheVersion)
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
