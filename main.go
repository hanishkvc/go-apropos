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
const PRG_VERSION = "v3-20220610IST1920"

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
var gbCaseSensitive bool

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

func handle_args() {
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
	flag.BoolVar(&gbCaseSensitive, "casesensitive", false, "Whether pkg name and symbol matching is case sensitive or not")
	flag.StringVar(&gsMatchMode, "matchmode", gsMatchMode, "Specify the strategy used for matching wrt pkg names and symbols. Supported modes contains regexp")
	flag.Parse()
	if len(flag.Args()) > 0 {
		if gFind != FIND_DUMMY {
			fmt.Printf("%v:WARN:ARG: Unknown args: %v\n", PRG_TAG, flag.Args())
			flag.Usage()
			os.Exit(1)
		}
		gFind = flag.Arg(0)
	}
	if (gFind == FIND_DUMMY) && (gFindPkg == FINDPKG_DEFAULT) && (gFindCmt == FINDCMT_DUMMY) {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "Sample usage:\n")
		fmt.Fprintln(os.Stderr, "\tgoapropos searchToken")
		fmt.Fprintln(os.Stderr, "\tgoapropos --find searchToken")
		fmt.Fprintln(os.Stderr, "\tgoapropos --findcmt searchToken")
		fmt.Fprintln(os.Stderr, "\tgoapropos --findpkg searchToken")
		fmt.Fprintln(os.Stderr, "\tgoapropos --matchmode regexp searchToken")
		fmt.Fprintln(os.Stderr, "Look at the README for more info about the program and its usage")
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
	}
}

func main() {
	handle_args()
	if giDEBUG > 1 {
		fmt.Println(PRG_TAG, PRG_NAME, PRG_VERSION)
	}
	test_go()
	do_walkdir(gBasePath)
	if giDEBUG > 3 {
		db_print()
	}
	if gFindPkg != FINDPKG_DEFAULT {
		db_print_pkgs()
	}
	db_find(gFind, gFindCmt)
}
