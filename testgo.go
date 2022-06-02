package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func test_data() {
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
	fmt.Println("After assigning slice and modifying assigned slice")
	fmt.Printf("assignedSlice: %v\n", assignedSlice)
	fmt.Printf("anIntSlice: %v\n", anIntSlice)
	// arrays
	var aArray [3]int = [3]int{10, 20, 30}
	var assignedArray [3]int
	fmt.Printf("aArray: %v\n", aArray)
	assignedArray = aArray
	assignedArray[1] = 888
	fmt.Println("After assigning array and modifying assigned array")
	fmt.Printf("aArray: %v\n", aArray)
	fmt.Printf("assignedArray: %v\n", assignedArray)
	// Maps
	var aMap map[string]int = map[string]int{"1": 1, "2": 2}
	fmt.Printf("aMap: %v\n", aMap)
	var bMap = map[string]int{"1": 1, "2": 2}
	fmt.Printf("bMap: %v\n", bMap)
	cMap := map[string]int{"1": 1, "2": 2}
	fmt.Printf("cMap: %v\n", cMap)
	var anotherMap = aMap
	anotherMap["2"] = 22
	fmt.Println("After assigning the map and modifying assigned map")
	fmt.Printf("anotherMap: %v\n", anotherMap)
	fmt.Printf("aMap: %v\n", aMap)
}

func test_flag() {
	piTest := flag.Int("int", 123, "Test a int flag")
	fmt.Printf("piTest: %v\n", piTest)
}

func test_fileread() {
	oF, err := os.Open("/etc/passwd")
	if err != nil {
		fmt.Printf("Open:err: %v\n", err)
		return
	}
	var buf [512]byte
	for {
		iCnt, err := oF.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nRead:Done:EOFRead")
			} else {
				fmt.Printf("\nRead:err: %v\n", err)
			}
			return
		}
		if iCnt == 0 {
			fmt.Println("\nRead:Done:0Read")
			break
		}
		fmt.Print(buf)
		fmt.Println("\nRead:INFO:", len(buf))
	}
}

func test_go() {
	fmt.Printf("%v:INFO: TestGo\n", PRG_TAG)
	test_data()
	test_fileread()
}
