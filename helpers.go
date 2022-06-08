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

func map_print(theMap any, sSep, sEnd string) {
	switch m := theMap.(type) {
	case map[string][]string:
		keys := []string{}
		for k, _ := range m {
			keys = append(keys, k)
		}
		keys = string_sort(keys)
		for _, k := range keys {
			fmt.Printf("%v%v%v%v", k, sSep, m[k], sEnd)
		}
	}
}
