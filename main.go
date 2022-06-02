package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
)

const PRG_TAG = "GOAPRO"
const PRG_NAME = "GoApropos"
const PRG_VERSION = "v0-20220602IST0954"

var gFind string
var gBasePath string = "/usr/lib/go-1.18/"

func handle_args() {
	flag.StringVar(&gFind, "find", "", "Specify the word to find")
	flag.StringVar(&gBasePath, "basepath", gBasePath, "Specify the dir containing files to search")
	flag.Parse()
	fmt.Printf("gFind: %v\n", gFind)
	fmt.Printf("gBasePath: %v\n", gBasePath)
	if len(flag.Args()) > 0 {
		fmt.Printf("%v:WARN: Unknown args: %v\n", PRG_TAG, flag.Args())
	}
}

func test_walkdir(sPath string) {
	oFS := os.DirFS(sPath)
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
		fmt.Printf("%v:INFO: %v:path: %v\n", PRG_TAG, sPType, path)
		return nil
	})
}

func main() {
	fmt.Println(PRG_NAME, PRG_VERSION)
	test_flag()
	handle_args()
	test_go()
	test_walkdir(gBasePath)
}
