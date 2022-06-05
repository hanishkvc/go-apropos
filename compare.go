// logic to compare tokens
// HanishKVC, 2022

package main

import "strings"

func match_prepare(sToken string) string {
	return strings.ToUpper(sToken)
}

func match_contains(sToCheck, sMatchTokenP string) bool {
	sToCheckP := strings.ToUpper(sToCheck)
	return strings.Contains(sToCheckP, sMatchTokenP)
}

func match_ok(sToCheck, sMatchTokenP string) bool {
	return match_contains(sToCheck, sMatchTokenP)
}
