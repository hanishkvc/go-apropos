// Provide equivalent of apropos wrt go packages
// HanishKVC, 2022

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const PRG_TAG = "GOAPRO"
const PRG_NAME = "GoApropos"
const PRG_VERSION = "v08-20220622IST2318"

const FIND_DUMMY = "__FIND_DUMMY__"
const FINDPKG_DEFAULT = ""
const FINDCMT_DUMMY = FIND_DUMMY

const BASEPATH_DEFAULT = "/usr/share/go-dummy/src"

var gFind string = FIND_DUMMY
var gFindPkg string = FINDPKG_DEFAULT
var gFindPkgP Matcher
var gFindCmt string = FINDCMT_DUMMY
var gBasePath string = BASEPATH_DEFAULT
var giDEBUG int = 0
var gbTEST bool = false
var gbAllSymbols bool = false
var gSkipFiles = []string{}
var gbCaseSensitive bool = false
var gbUseCache bool = false
var gbCreateCache bool = false
var gbIndentJSON bool = false
var gbSortedResult bool = false
var gbAutoCache bool = true

var giMatchMode = MatchMode_RegExp
var gsMatchMode string = MATCHMODE_REGEXP

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
		sPath := filepath.Join(basePath, sDirName, srcDir)
		srcPaths = append(srcPaths, sPath)
	}
	if giDEBUG > 0 { // Needs to be enabled by setting giDEBUG in source
		fmt.Printf("%v:INFO:FindSrcPaths: basePath: %v, srcPaths: %v\n", PRG_TAG, basePath, srcPaths)
	}
	return srcPaths
}

// As the path at which the go src files are stored can change btw
// systems, so build a list of possible directories based on GOROOT
// and some builtin predefined possibilities and then select the
// 1st path which seems to exist.
//
// This should potentially allow the logic to work across atleast few
// different OSs and Distros with their own path wrt go installation.
func set_gbasepath() {
	srcPaths := []string{filepath.Join(runtime.GOROOT(), "src")}
	for _, lookAt := range []string{"/usr/share", "/usr/local/share", "/usr/lib"} {
		srcPaths = find_srcpaths(lookAt, srcPaths)
	}
	for _, srcPath := range srcPaths {
		fInfo, err := os.Stat(srcPath)
		if err != nil {
			if giDEBUG > 0 {
				fmt.Printf("%v:WARN:SetGBasePath: %v, %v\n", PRG_TAG, srcPath, err)
			}
			continue
		}
		if giDEBUG > 0 {
			fmt.Printf("%v:INFO:SetGBasePath: %v [N:%v, D:%v, M:%v, T:%v, S:%v]\n", PRG_TAG, srcPath, fInfo.Name(), fInfo.IsDir(), fInfo.Mode(), fInfo.ModTime(), fInfo.Size())
		}
		gBasePath = srcPath
		return
	}
}

var sAdditional string = `
Sample usage:

	goapropos searchToken
	goapropos --find searchToken
		The above two are equivalent and search for a matching symbol across all the packages

	goapropos --findcmt searchToken
		This searchs thro the comments wrt all the packages and their symbols for a match.

	goapropos --findpkg searchToken
		Find packages whose name match the given search token

	goapropos --findpkg pkgNameSearchToken --find symbolSearchToken
		Find symbols which match the symbolSearchToken, from across all packages whose name match pkgNameSearchToken

	goapropos --matchmode contains searchToken
		Use contains-substring matching logic, instead of the regexp based default logic.

	goapropos --autocache=false --createcache
		Force the recreation of the program's internal cache.

	goapropos --autocache=false --findpkg pkgNameSearchToken symbolSearchToken
		This disables internal cache and parses through all go source files to find any matching stuff.

Look at the README for more info about the program and its usage`

func handle_args() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%v:%v:%v\n", PRG_TAG, PRG_NAME, PRG_VERSION)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), sAdditional)
	}

	set_gbasepath()
	flag.StringVar(&gFind, "find", gFind, "Specify the pattern/substring to match wrt symbols. The pattern to match can also be specified as a standalone arg on its own")
	flag.StringVar(&gFindPkg, "findpkg", gFindPkg, "Specify the pattern/substring to match wrt package name")
	flag.StringVar(&gFindCmt, "findcmt", gFindCmt, "Specify the pattern/substring to match wrt comments in package source")
	flag.StringVar(&gBasePath, "basepath", gBasePath, "Specify the go src dir containing src files to search")
	flag.IntVar(&giDEBUG, "debug", 0, "Set debug level to control debug prints")
	flag.BoolVar(&gbTEST, "test", gbTEST, "Enable test logics")
	flag.BoolVar(&gbAllSymbols, "allsymbols", gbAllSymbols, "Match all symbols and not just exported. [NEEDs:createcache OR autocache=false]")
	flag.Func("skipfiles", "Specify pattern to match wrt package path+filename for skipping package files. More than one can be specified. [NEEDs: createcache OR autocache=False]", func(s string) error {
		gSkipFiles = append(gSkipFiles, s)
		return nil
	})
	flag.BoolVar(&gbCaseSensitive, "casesensitive", gbCaseSensitive, "Whether pkg name and symbol matching is case sensitive or not")
	flag.StringVar(&gsMatchMode, "matchmode", gsMatchMode, "Specify the strategy used for matching. Supported modes are contains regexp")
	flag.BoolVar(&gbAutoCache, "autocache", gbAutoCache, "Create and use the cache automatically. This manipulates usecache and createcache flags automatically")
	flag.BoolVar(&gbCreateCache, "createcache", gbCreateCache, "Create a cache of the package symbols, paths and comments")
	flag.BoolVar(&gbUseCache, "usecache", gbUseCache, "Use cache of the package symbols, paths and comments, instead of parsing the go sources")
	flag.BoolVar(&gbIndentJSON, "indentjson", gbIndentJSON, "Create pretty indented json cache files")
	flag.BoolVar(&gbSortedResult, "sortedresult", gbSortedResult, "Show results as found or sorted at the end")
	flag.Parse()
	if gBasePath == BASEPATH_DEFAULT {
		fmt.Fprintf(flag.CommandLine.Output(), "%v:ERRR:SetGBasePath: No valid Go src path found, please specify using --basepath", PRG_TAG)
		os.Exit(1)
	}
	cache_maya()
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
	giMatchMode = matchmode_fromstr(gsMatchMode)
	gFindPkgP = New_Matcher(giMatchMode, gFindPkg, gbCaseSensitive)
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
	bWalkPlus := true
	if gbUseCache {
		err := load_dbs()
		if err != nil {
			fmt.Printf("%v:ERRR:Main: flying failed, will walk and cache...[%v]\n", PRG_TAG, err)
			gbCreateCache = true
		} else {
			bWalkPlus = false
		}
	}
	if bWalkPlus {
		do_walkdir(gBasePath)
		if gbCreateCache {
			prep_dir()
			save_dbs()
		}
	}
	if giDEBUG > 3 {
		dbprint_all(gDB, "Package:", "\n", "\t", "\n", "\n")
	}
	db_find(gDB, gFind, gFindCmt, gFindPkg, giMatchMode, gbCaseSensitive, gbSortedResult)
}
