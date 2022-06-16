// Provide equivalent of apropos wrt go packages
// HanishKVC, 2022

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const PRG_TAG = "GOAPRO"
const PRG_NAME = "GoApropos"
const PRG_VERSION = "v6-20220616IST1710"

const FIND_DUMMY = "__FIND_DUMMY__"
const FINDPKG_DEFAULT = ""
const FINDCMT_DUMMY = FIND_DUMMY

var gFind string = FIND_DUMMY
var gFindPkg string = FINDPKG_DEFAULT
var gFindPkgP string = match_prepare(gFindPkg) // the explicit initialising can be avoided, but still
var gFindCmt string = FINDCMT_DUMMY
var gBasePath string = "/usr/share/go-dummy/"
var giDEBUG int = 0
var gbTEST bool
var gbAllSymbols bool
var gSkipFiles = []string{}
var gbCaseSensitive bool = false
var gbUseCache bool = false
var gbCreateCache bool = false
var gbIndentJSON bool = false
var gbSortedResult bool = false
var gbAutoCache bool = true

var giMatchMode = MatchMode_Contains
var gsMatchMode string = "contains"

func find_srcpaths(basePath string, srcPaths []string) []string {
	namePrefixs := []string{"go-", "golang"}
	const srcDir = "src"
	aDE, err := os.ReadDir(basePath)
	if err != nil {
		if giDEBUG > 0 { // Needs to be enabled by setting giDEBUG in source
			fmt.Printf("%v:ERRR:FindSrcPaths: basePath: %v, Err: %v\n", PRG_TAG, basePath, err)
		}
		return srcPaths
	}
	for _, de := range aDE {
		if !de.IsDir() {
			continue
		}
		sDirName := de.Name()
		if !string_hasprefix_anysubstring(sDirName, namePrefixs) {
			continue
		}
		sPath := strings.Join([]string{basePath, sDirName, srcDir}, string(os.PathSeparator))
		srcPaths = append(srcPaths, sPath)
	}
	if giDEBUG > 0 { // Needs to be enabled by setting giDEBUG in source
		fmt.Printf("%v:INFO:FindSrcPaths: basePath: %v, srcPaths: %v\n", PRG_TAG, basePath, srcPaths)
	}
	return srcPaths
}

func set_gbasepath() {
	srcPaths := []string{}
	for _, lookAt := range []string{"/usr/share", "/usr/local/share", "/usr/lib"} {
		srcPaths = find_srcpaths(lookAt, srcPaths)
	}
	if len(srcPaths) > 0 {
		gBasePath = srcPaths[0]
	}
}

var sAdditional string = `
Sample usage:
	goapropos searchToken
	goapropos --find searchToken
	goapropos --findcmt searchToken
	goapropos --findpkg searchToken
	goapropos --matchmode regexp searchToken
	goapropos --createcache
	goapropos --usecache --matchmode regexp --findcmt searchToken
Look at the README for more info about the program and its usage`

func handle_args() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%v:%v:%v\n", PRG_TAG, PRG_NAME, PRG_VERSION)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), sAdditional)
	}

	set_gbasepath()
	flag.StringVar(&gFind, "find", gFind, "Specify the token/substring to match wrt symbols. The token to match can also be specified as a standalone arg on its own")
	flag.StringVar(&gFindPkg, "findpkg", gFindPkg, "Specify the token/substring to match wrt package name")
	flag.StringVar(&gFindCmt, "findcmt", gFindCmt, "Specify the token/substring to match wrt comments in package source")
	flag.StringVar(&gBasePath, "basepath", gBasePath, "Specify the dir containing files to search")
	flag.IntVar(&giDEBUG, "debug", 0, "Set debug level to control debug prints")
	flag.BoolVar(&gbTEST, "test", false, "Enable test logics")
	flag.BoolVar(&gbAllSymbols, "allsymbols", false, "Match all symbols and not just exported")
	flag.Func("skipfiles", "Specify token to match wrt package path+filename for skipping package files. More than one can be specified", func(s string) error {
		gSkipFiles = append(gSkipFiles, s)
		return nil
	})
	flag.BoolVar(&gbCaseSensitive, "casesensitive", gbCaseSensitive, "Whether pkg name and symbol matching is case sensitive or not")
	flag.StringVar(&gsMatchMode, "matchmode", gsMatchMode, "Specify the strategy used for matching wrt pkg names and symbols. Supported modes contains regexp")
	flag.BoolVar(&gbAutoCache, "autocache", gbAutoCache, "Create and use the cache automatically. This manipulates usecache and createcache flags automatically")
	flag.BoolVar(&gbCreateCache, "createcache", gbCreateCache, "Create a cache of the package symbols, paths and comments")
	flag.BoolVar(&gbUseCache, "usecache", gbUseCache, "Use cache of the package symbols, paths and comments, instead of parsing the go sources")
	flag.BoolVar(&gbIndentJSON, "indentjson", gbIndentJSON, "Create pretty indented json cache files")
	flag.BoolVar(&gbSortedResult, "sortedresult", gbSortedResult, "Show results as found or sorted at the end")
	flag.Parse()
	if gbAutoCache {
		cache_maya()
	}
	if len(flag.Args()) > 0 {
		if gFind != FIND_DUMMY {
			fmt.Printf("%v:WARN:ARG: Unknown args: %v\n", PRG_TAG, flag.Args())
			flag.Usage()
			os.Exit(1)
		}
		gFind = flag.Arg(0)
	}
	if ((gFind == FIND_DUMMY) && (gFindPkg == FINDPKG_DEFAULT) && (gFindCmt == FINDCMT_DUMMY)) && !gbCreateCache {
		flag.Usage()
		os.Exit(1)
	}
	gFindPkgP = match_prepare(gFindPkg)
	giMatchMode = matchmode_fromstr(gsMatchMode)
	if giDEBUG > 1 {
		fmt.Printf("%v:INFO:ARG: gFind: %v\n", PRG_TAG, gFind)
		fmt.Printf("%v:INFO:ARG: gFindPkg: %v\n", PRG_TAG, gFindPkg)
		fmt.Printf("%v:INFO:ARG: gFindCmt: %v\n", PRG_TAG, gFindCmt)
		fmt.Printf("%v:INFO:ARG: gBasePath: %v\n", PRG_TAG, gBasePath)
		fmt.Printf("%v:INFO:ARG: giDEBUG: %v\n", PRG_TAG, giDEBUG)
		fmt.Printf("%v:INFO:ARG: gbAllSymbols: %v\n", PRG_TAG, gbAllSymbols)
		fmt.Printf("%v:INFO:ARG: gbTEST: %v\n", PRG_TAG, gbTEST)
		fmt.Printf("%v:INFO:ARG: gSkipFiles: %v\n", PRG_TAG, gSkipFiles)
		fmt.Printf("%v:INFO:ARG: gbCaseSensitive: %v\n", PRG_TAG, gbCaseSensitive)
		fmt.Printf("%v:INFO:ARG: gsMatchMode: %v\n", PRG_TAG, gsMatchMode)
		fmt.Printf("%v:INFO:ARG: gbCreateCache: %v\n", PRG_TAG, gbCreateCache)
		fmt.Printf("%v:INFO:ARG: gbUseCache: %v\n", PRG_TAG, gbUseCache)
		fmt.Printf("%v:INFO:ARG: gbIndentJSON: %v\n", PRG_TAG, gbIndentJSON)
		fmt.Printf("%v:INFO:ARG: gbSortedResult: %v\n", PRG_TAG, gbSortedResult)
	}
}

func main() {
	handle_args()
	if giDEBUG > 1 {
		fmt.Println(PRG_TAG, PRG_NAME, PRG_VERSION)
	}
	test_go()
	if gbUseCache {
		load_dbs()
	} else {
		do_walkdir(gBasePath)
		if gbCreateCache {
			prep_dir()
			save_dbs()
		}
	}
	if giDEBUG > 3 {
		dbprint_all_all(gDB)
	}
	db_find(gDB, gFind, gFindCmt, gFindPkg)
}
