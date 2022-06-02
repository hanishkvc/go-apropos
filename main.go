package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
)

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
}

func test_walkdir(sPath string) {
	oFS := os.DirFS(sPath)
	fs.WalkDir(oFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("GOAP:ERRR: path: %v\n", path)
			return err
		}
		fmt.Printf("GOAP:INFO: path: %v, err: %v\n", path, err)
		if d.Type().IsDir() {
			fmt.Println("GOAP:INFO: dir")
		}
		return err
	})
}

func main() {
	fmt.Println(PRG_NAME, PRG_VERSION)
	handle_args()
	test_walkdir(gBasePath)
}
