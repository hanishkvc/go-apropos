package main

import "fmt"

func test_go() {
	fmt.Printf("%v:INFO: TestGo\n", PRG_TAG)
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
}
