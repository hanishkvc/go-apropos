// A set of helper routines
// HanishKVC, 2022

package main

import (
	"fmt"
	"sort"
	"strings"
)

func substring_inthere(theStrings []string, subString string) bool {
	for _, curString := range theStrings {
		if strings.Contains(curString, subString) {
			return true
		}
	}
	return false
}

func string_contains_anysubstring(theString string, subStrings []string) bool {
	for _, curSubString := range subStrings {
		if strings.Contains(theString, curSubString) {
			return true
		}
	}
	return false
}

func string_hasprefix_anysubstring(theString string, subStrings []string) bool {
	for _, curSubString := range subStrings {
		if strings.HasPrefix(theString, curSubString) {
			return true
		}
	}
	return false
}

func string_sort(theSlice []string) []string {
	sort.SliceStable(theSlice, func(i, j int) bool {
		return theSlice[i] < theSlice[j]
	})
	return theSlice
}

// THis is a sort of centralised map print logic
// using pure Interface type mechanism to create a sort of generic function.
func map_print(theMap map[string]any, sPrefix, sPrefixSep, sSep, sEnd string) {
	keys := sort.StringSlice{}
	for k, _ := range theMap {
		keys = append(keys, k)
	}
	keys.Sort()
	for _, k := range keys {
		v := theMap[k]
		switch t := v.(type) {
		case []string:
			sort.Strings(t)
		}
		fmt.Printf("%v%v%v%v%v%v", sPrefix, sPrefixSep, k, sSep, v, sEnd)
	}
}

// THis is a sort of generic map print logic
// using Type parameters support in the latest Go language versions (>= 1.18)
func map_print_go118[GT any](theMap map[string]GT, sPrefix, sPrefixSep, sSep, sEnd string) {
	keys := sort.StringSlice{}
	for k, _ := range theMap {
		keys = append(keys, k)
	}
	keys.Sort()
	for _, k := range keys {
		fmt.Printf("%v%v%v%v%v%v", sPrefix, sPrefixSep, k, sSep, theMap[k], sEnd)
	}
}
