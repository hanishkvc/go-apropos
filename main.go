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
var gBasePath string = "/usr/lib/go-1.18/"
var giDEBUG int
var gbTEST bool
var gbALL bool

func handle_args() {
	flag.StringVar(&gFind, "find", "", "Specify the word to find")
	flag.StringVar(&gBasePath, "basepath", gBasePath, "Specify the dir containing files to search")
	flag.IntVar(&giDEBUG, "debug", 0, "Set debug level to control debug prints")
	flag.BoolVar(&gbTEST, "test", false, "Enable test logics")
	flag.BoolVar(&gbALL, "all", false, "Match all symbols and not just exported")
	flag.Parse()
	if giDEBUG > 0 {
		fmt.Printf("gFind: %v\n", gFind)
		fmt.Printf("gBasePath: %v\n", gBasePath)
		fmt.Printf("giDEBUG: %v\n", giDEBUG)
		fmt.Printf("gbALL: %v\n", gbALL)
		fmt.Printf("gbTEST: %v\n", gbTEST)
	}
	if len(flag.Args()) > 0 {
		fmt.Printf("%v:WARN: Unknown args: %v\n", PRG_TAG, flag.Args())
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
	if giDEBUG > 5 {
		fmt.Printf("%v:INFO: GoPkg:%v:%v\n", PRG_TAG, name, idents)
	}
}

func do_walkdir(sPath string) {
	oFS := os.DirFS(sPath)
	if giDEBUG > 10 {
		fmt.Printf("oFS: %v\n", oFS)
	}
	fs.WalkDir(oFS, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%v:ERRR: path: %v, Err:%v\n", PRG_TAG, path, err)
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
		if giDEBUG > 1 {
			fmt.Printf("%v:INFO: %v:path: %v\n", PRG_TAG, sPType, path)
		}
		if sPType == "File" {
			theFile := sPath + "/" + path
			handle_file(theFile)
		}
		return nil
	})
}

func main() {
	if giDEBUG > 0 {
		fmt.Println(PRG_TAG, PRG_NAME, PRG_VERSION)
	}
	test_flag()
	handle_args()
	test_go()
	do_walkdir(gBasePath)
	if giDEBUG > 2 {
		db_print()
	}
	db_find(gFind)
}
