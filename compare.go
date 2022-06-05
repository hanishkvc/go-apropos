// logic to compare tokens
// HanishKVC, 2022

package main

import (
	"regexp"
	"strings"
)

const MATCHMODE_REGEXP = "regexp"
const MATCHMODE_CONTAINS = "contains"

// Prepare a token / string for use by match_ok logic
func match_prepare(sToken string) string {
	if gbCaseSensitive {
		return sToken
	}
	return strings.ToUpper(sToken)
}

// Check if sToCheck contains sMatchTokenP as a substring with in it or not.
//
// sToCheck is expected to be the raw string.
// sMatchTokenP is expected to be the token string, which has already been processed/prepared using match_prepare.
func match_contains(sToCheck, sMatchTokenP string) bool {
	sToCheckP := match_prepare(sToCheck)
	return strings.Contains(sToCheckP, sMatchTokenP)
}

func match_regexp(sToCheck, sMatchTokenP string) bool {
	sToCheckP := match_prepare(sToCheck)
	match, err := regexp.Match(sMatchTokenP, []byte(sToCheckP))
	if err != nil {
		return false
	}
	return match
}

// Check if sToCheck satisfying the match token sMatchTokenP or not.
//
// It might use different strategies to check for a match like contains or regexp or ...
//
// sToCheck is expected to be the raw string.
// sMatchTokenP is expected to be the token string, which has already been processed/prepared using match_prepare.
func match_ok(sToCheck, sMatchTokenP string) bool {
	if gsMatchMode == MATCHMODE_REGEXP {
		return match_regexp(sToCheck, sMatchTokenP)
	}
	return match_contains(sToCheck, sMatchTokenP)
}
