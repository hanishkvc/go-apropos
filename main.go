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

func test_go() {
	var anIntSlice []int = []int{1, 2, 3}
	var assignedSlice []int
	anIntSlice = append(anIntSlice, 4, 5, 6)
	anIntSlice = append(anIntSlice, []int{10, 20, 30}...)
	fmt.Printf("anIntArray: %v\n", anIntSlice)
	fmt.Printf("len(anIntSlice): %v\n", len(anIntSlice))
	var anotherSlice []int
	fmt.Printf("anotherSlice: %v\n", anotherSlice)
	fmt.Printf("len(anotherSlice): %v\n", len(anotherSlice))
	assignedSlice = anIntSlice
	assignedSlice[1] = 999
	fmt.Printf("assignedSlice: %v\n", assignedSlice)
	fmt.Printf("anIntSlice: %v\n", anIntSlice)
	// arrays
	var aArray [3]int = [3]int{10, 20, 30}
	var assignedArray [3]int
	fmt.Printf("aArray: %v\n", aArray)
	assignedArray = aArray
	assignedArray[1] = 888
	fmt.Printf("aArray: %v\n", aArray)
	fmt.Printf("assignedArray: %v\n", assignedArray)
}

func handle_args() {
	piTest := flag.Int("int", 123, "Test a int flag")
	fmt.Printf("piTest: %v\n", piTest)
	flag.StringVar(&gFind, "find", "", "Specify the word to find")
	flag.StringVar(&gBasePath, "basepath", gBasePath, "Specify the dir containing files to search")
	flag.Parse()
	fmt.Printf("gFind: %v\n", gFind)
	fmt.Printf("gBasePath: %v\n", gBasePath)
	fmt.Printf("%v:WARN: Unknown args: %v\n", PRG_TAG, flag.Args())
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
	handle_args()
	test_go()
	test_walkdir(gBasePath)
}
