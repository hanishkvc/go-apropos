// A set of helper routines
// HanishKVC, 2022

package main

import (
	"fmt"
	"sort"
)

func string_sort(theSlice []string) []string {
	sort.SliceStable(theSlice, func(i, j int) bool {
		return theSlice[i] < theSlice[j]
	})
	return theSlice
}

// THis is a sort of centralised map print logic
// using pure Interface type mechanism to create a sort of generic function.
func map_print(theMap any, sSep, sEnd string) {
	switch m := theMap.(type) {
	case map[string][]string:
		keys := sort.StringSlice{}
		for k, _ := range m {
			keys = append(keys, k)
		}
		keys.Sort()
		for _, k := range keys {
			fmt.Printf("%v%v%v%v", k, sSep, m[k], sEnd)
		}
	}
}

// THis is a sort of generic map print logic
// using Type parameters support in the latest Go language versions (>= 1.18)
func map_print_go118[GT any](theMap map[string]GT, sSep, sEnd string) {
	keys := sort.StringSlice{}
	for k, _ := range theMap {
		keys = append(keys, k)
	}
	keys.Sort()
	for _, k := range keys {
		fmt.Printf("%v%v%v%v", k, sSep, theMap[k], sEnd)
	}
}
