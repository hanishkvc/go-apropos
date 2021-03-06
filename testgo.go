// Some code to understand Go
// HanishKVC, 2022

package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

const FILE2READ = "/etc/passwd"

func test_data() {
	var anIntSlice []int = []int{1, 2, 3}
	var assignedSlice []int
	anIntSlice = append(anIntSlice, 4, 5, 6)
	anIntSlice = append(anIntSlice, []int{10, 20, 30}...)
	fmt.Printf("%v:INFO:T SLICE: anIntSlice: %v\n", PRG_TAG, anIntSlice)
	fmt.Printf("%v:INFO:T SLICE: len(anIntSlice): %v\n", PRG_TAG, len(anIntSlice))
	var anotherSlice []int
	fmt.Printf("%v:INFO:T SLICE: anotherSlice: %v\n", PRG_TAG, anotherSlice)
	fmt.Printf("%v:INFO:T SLICE: len(anotherSlice): %v\n", PRG_TAG, len(anotherSlice))
	assignedSlice = anIntSlice
	assignedSlice[1] = 999
	fmt.Printf("%v:INFO:T SLICE: After assigning slice and modifying assigned slice\n", PRG_TAG)
	fmt.Printf("%v:INFO:T SLICE: assignedSlice: %v\n", PRG_TAG, assignedSlice)
	fmt.Printf("%v:INFO:T SLICE: anIntSlice: %v\n", PRG_TAG, anIntSlice)
	// arrays
	var aArray [3]int = [3]int{10, 20, 30}
	var assignedArray [3]int
	fmt.Printf("%v:INFO:T ARRAY: aArray: %v\n", PRG_TAG, aArray)
	assignedArray = aArray
	assignedArray[1] = 888
	fmt.Printf("%v:INFO:T ARRAY: After assigning array and modifying assigned array\n", PRG_TAG)
	fmt.Printf("%v:INFO:T ARRAY: aArray: %v\n", PRG_TAG, aArray)
	fmt.Printf("%v:INFO:T ARRAY: assignedArray: %v\n", PRG_TAG, assignedArray)
	// struct
	var aStruct struct {
		cnt int
		doc string
	}
	type TestStruct struct {
		cnt int
		doc string
	}
	aStruct.doc = "test doc value"
	var assignedStruct TestStruct // Initialisation below kept seperate from var declaration
	assignedStruct = aStruct
	fmt.Printf("%v:INFO:T STRUCT: aStruct: %v\n", PRG_TAG, aStruct)
	aStruct.cnt += 1
	assignedStruct.cnt = 99
	fmt.Printf("%v:INFO:T STRUCT: After assigning struct and modifying assigned struct; orig struct part update also\n", PRG_TAG)
	fmt.Printf("%v:INFO:T STRUCT: aStruct: %v\n", PRG_TAG, aStruct)
	fmt.Printf("%v:INFO:T STRUCT: assignedStruct: %v\n", PRG_TAG, assignedStruct)
	// Maps
	var aMap map[string]int = map[string]int{"1": 1, "2": 2}
	fmt.Printf("%v:INFO:T MAP: aMap: %v\n", PRG_TAG, aMap)
	var bMap = map[string]int{"1": 1, "2": 2}
	fmt.Printf("%v:INFO:T MAP: bMap: %v\n", PRG_TAG, bMap)
	cMap := map[string]int{"1": 1, "2": 2}
	fmt.Printf("%v:INFO:T MAP: cMap: %v\n", PRG_TAG, cMap)
	var assignedMap = aMap
	aMap["1"] += 1
	assignedMap["2"] = 22
	fmt.Printf("%v:INFO:T MAP: After assigning the map and modifying assigned map; orig map part update also\n", PRG_TAG)
	fmt.Printf("%v:INFO:T MAP: assignedMap: %v\n", PRG_TAG, assignedMap)
	fmt.Printf("%v:INFO:T MAP: aMap: %v\n", PRG_TAG, aMap)
	// Map with Struct Value
	var aMapS = map[string]TestStruct{}
	// aMapS["t1"] = Ident{0, "test1"} // Not allowed as Named Structs ie Type names are different even though internal structure same
	aMapS["t1"] = TestStruct{0, "test1"}
	fmt.Printf("%v:INFO:T MAPSTRUCT: aMapS: %v\n", PRG_TAG, aMapS)
	fmt.Printf("%v:INFO:T MAPSTRUCT: aMapS[\"t1\"].cnt: %v\n", PRG_TAG, aMapS["t1"].cnt)
	//aMapS["t1"].cnt = 99
	//aMapS["t1"].cnt += 1
	tVS, _ := aMapS["t1"]
	//tVS.cnt += 1
	tVS.cnt = 99
	fmt.Printf("%v:INFO:T MAPSTRUCT: After assigning struct value of map and modifying assigned struct\n", PRG_TAG)
	fmt.Printf("%v:INFO:T MAPSTRUCT: aMapS: %v\n", PRG_TAG, aMapS)
	fmt.Printf("%v:INFO:T MAPSTRUCT: tVS: %v\n", PRG_TAG, tVS)
	// Map with Struct Value Test2
	test_dbsymbols_update()
}

func test_assignment() {
	fmt.Printf("%v:INFO:T ASSIGNMENTS: Bunch of value assignments follow\n", PRG_TAG)

	type MyInt int
	var aInt int
	var mInt MyInt = 123
	aInt = int(mInt)
	fmt.Printf("\t %T:%v, %T:%v\n", aInt, aInt, mInt, mInt)

	type MyRune rune
	var aRune rune
	var mRune MyRune = 234
	aRune = rune(mRune)
	fmt.Printf("\t %T:%v, %T:%v\n", aRune, aRune, mRune, mRune)

	type MyString string
	var aString string
	var mString MyString = ""
	aString = string(mString)
	fmt.Printf("\t %T:%v, %T:%v\n", aString, aString, mString, mString)

	type MyIntArray [3]int
	var aIntArray [3]int
	var mIntArray MyIntArray
	aIntArray = mIntArray
	fmt.Printf("\t %T:%v, %T:%v\n", aIntArray, aIntArray, mIntArray, mIntArray)

	type MyIntSlice []int
	var aIntSlice []int
	var mIntSlice MyIntSlice
	aIntSlice = mIntSlice
	fmt.Printf("\t %T:%v, %T:%v\n", aIntSlice, aIntSlice, mIntSlice, mIntSlice)

	type MyStruct struct {
		f1 int
		f2 string
	}
	var aStruct = struct {
		f1 int
		f2 string
	}{11, "11"}
	var mStruct MyStruct
	aStruct = mStruct
	fmt.Printf("\t %T:%v, %T:%v\n", aStruct, aStruct, mStruct, mStruct)

	type MyMap map[string]int
	var aMap map[string]int
	var mMap MyMap
	aMap = mMap
	fmt.Printf("\t %T:%v, %T:%v\n", aMap, aMap, mMap, mMap)

	xMap := map[string]DBSymbolInfo{"t100": {"doct of t100", "X"}}
	yMap := DBSymbols{"t200": {"doc of t200", "Y"}}
	xMap = yMap
	fmt.Printf("\t %T:%v, %T:%v\n", xMap, xMap, yMap, yMap)

}

func test_gochan() {
	c1 := make(chan int, 3)
	go func() {
		for {
			x := <-c1
			fmt.Printf("%v:INFO:T GOCHAN: Got:%v\n", PRG_TAG, x)
			time.Sleep(2 * time.Second)
		}
	}()
	for i := 0; i < 10; i++ {
		t1 := time.Now()
		c1 <- i
		fmt.Printf("%v:INFO:T GOCHAN: Send:%v:Took:%v\n", PRG_TAG, i, time.Since(t1))
	}
}

func test_flag() {
	var lInt int
	piTest := flag.Int("int", 123, "Test a int flag")
	flag.IntVar(&lInt, "localint", 999, "test a local to function int flag")
	fmt.Printf("%v:INFO:T FLAG: piTest: %v\n", PRG_TAG, piTest)
	fmt.Printf("%v:INFO:T FLAG: &lInt: %v\n", PRG_TAG, &lInt)
	fmt.Printf("%v:INFO:T FLAG: &gbTEST: %v\n", PRG_TAG, &gbTEST)
	fmt.Printf("%v:INFO:T FLAG: &gFind: %v\n", PRG_TAG, &gFind)
}

func test_fileread_low(sFilePath string) {
	fmt.Printf("%v:INFO: TestFileRead:Low\n", PRG_TAG)
	oF, err := os.Open(sFilePath)
	if err != nil {
		fmt.Printf("%v:ERRR:T FREADLOW: Open: %v\n", PRG_TAG, err)
		return
	}
	var buf [512]byte
	for {
		iCnt, err := oF.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				fmt.Printf("\n%v:INFO:T FREADLOW: Done:EOFRead\n", PRG_TAG)
			} else {
				fmt.Printf("\n%v:ERRR:T FREADLOW: Read: %v\n", PRG_TAG, err)
			}
			return
		}
		if iCnt == 0 {
			fmt.Printf("\n%v:INFO:T FREADLOW: Done:0Read\n", PRG_TAG)
			break
		}
		fmt.Print(string(buf[:iCnt]))
		if giDEBUG > 10 {
			fmt.Printf("\n%v:INFO:T FREADLOW: ReadAmount: BufLen:%v ReadRet:%v\n", PRG_TAG, len(buf), iCnt)
		}
	}
}

func test_fileread_simple(sFilePath string) {
	fmt.Printf("%v:INFO: TestFileRead:Simple\n", PRG_TAG)
	oFS := os.DirFS("/")
	bData, err := fs.ReadFile(oFS, sFilePath)
	if err != nil {
		fmt.Printf("%v:ERRR:T FREADSIMP: %v\n", PRG_TAG, err)
		return
	}
	fmt.Printf("%v:INFO:T FREADSIMP: Read:\n%v\n", PRG_TAG, string(bData))
}

func test_env() {
	//fmt.Printf("%v:INFO:T ENVIRON: os.Environ(): %v\n", PRG_TAG, os.Environ())
	for _, env := range os.Environ() {
		fmt.Printf("%v:INFO:T ENVIRON: %v\n", PRG_TAG, env)
	}
}

func test_go() {
	if !gbTEST {
		return
	}
	fmt.Printf("%v:INFO: TestGo\n", PRG_TAG)
	test_flag()
	test_data()
	test_assignment()
	test_env()
	test_if()
	test_gochan()
	test_fileread_low(FILE2READ)
	test_fileread_simple(FILE2READ[1:])
}
