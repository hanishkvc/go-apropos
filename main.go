package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const PRG_TAG = "GOAPRO"
const PRG_NAME = "GoApropos"
const PRG_VERSION = "v1-20220604IST0942"

var gFind string
var gBasePath string = "/usr/share/go-dummy/"
var giDEBUG int = 0
var gbTEST bool
var gbALL bool

func find_srcpaths(basePath string, srcPaths []string) []string {
	const namePrefix = "go-"
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
		if !strings.HasPrefix(sDirName, namePrefix) {
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
	for _, lookAt := range []string{"/usr/share", "/usr/local/share"} {
		srcPaths = find_srcpaths(lookAt, srcPaths)
	}
	if len(srcPaths) > 0 {
		gBasePath = srcPaths[0]
	}
}

func handle_args() {
	set_gbasepath()
	flag.StringVar(&gFind, "find", "", "Specify the word to find")
	flag.StringVar(&gBasePath, "basepath", gBasePath, "Specify the dir containing files to search")
	flag.IntVar(&giDEBUG, "debug", 0, "Set debug level to control debug prints")
	flag.BoolVar(&gbTEST, "test", false, "Enable test logics")
	flag.BoolVar(&gbALL, "all", false, "Match all symbols and not just exported")
	flag.Parse()
	if giDEBUG > 1 {
		fmt.Printf("%v:INFO:ARG: gFind: %v\n", PRG_TAG, gFind)
		fmt.Printf("%v:INFO:ARG: gBasePath: %v\n", PRG_TAG, gBasePath)
		fmt.Printf("%v:INFO:ARG: giDEBUG: %v\n", PRG_TAG, giDEBUG)
		fmt.Printf("%v:INFO:ARG: gbALL: %v\n", PRG_TAG, gbALL)
		fmt.Printf("%v:INFO:ARG: gbTEST: %v\n", PRG_TAG, gbTEST)
	}
	if len(flag.Args()) > 0 {
		fmt.Printf("%v:WARN:ARG: Unknown args: %v\n", PRG_TAG, flag.Args())
	}
}

func handle_file(sFile string) {
	if !strings.HasSuffix(sFile, "go") {
		return
	}
	if strings.HasSuffix(sFile, "_test.go") {
		return
	}
	name, idents := gosrc_info(sFile)
	db_add(name, idents)
}

func do_walkdir(sPath string) {
	oFS := os.DirFS(sPath)
	if giDEBUG > 10 {
		fmt.Printf("%v:INFO:WALKDIR: oFS: %v\n", PRG_TAG, oFS)
	}
	fs.WalkDir(oFS, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%v:ERRR:WALKDIR: path: %v, Err: %v\n", PRG_TAG, path, err)
			return err
		}
		var sPType string
		deT := de.Type()
		if deT.IsDir() {
			sPType = "Dir"
		} else if deT.IsRegular() {
			sPType = "File"
		} else {
			sPType = "???"
		}
		if giDEBUG > 2 {
			fmt.Printf("%v:INFO:WALKDIR: %v:path: %v\n", PRG_TAG, sPType, path)
		}
		if sPType == "File" {
			theFile := sPath + string(os.PathSeparator) + path
			handle_file(theFile)
		}
		return nil
	})
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
	db_find(gFind)
}
