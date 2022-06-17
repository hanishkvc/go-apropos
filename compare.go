// logic to compare tokens
// HanishKVC, 2022

package main

import (
	"fmt"
	"regexp"
	"strings"
)

const MATCHMODE_CONTAINS = "contains"
const MATCHMODE_REGEXP = "regexp"

type MatchMode uint

const (
	MatchMode_Contains MatchMode = iota
	MatchMode_RegExp
)

func matchmode_fromstr(sMode string) MatchMode {
	switch sMode {
	case MATCHMODE_CONTAINS:
		return MatchMode_Contains
	case MATCHMODE_REGEXP:
		return MatchMode_RegExp
	default:
		return MatchMode_Contains
	}
}

type UMTP_string string
type UMTP_re regexp.Regexp
type UMTP interface {
	Utype() string
	Matchok(string) bool
}

func (o UMTP_string) Utype() string {
	return "string"
}

func (subStr UMTP_string) Matchok(theStr string) bool {
	return strings.Contains(theStr, string(subStr))
}

func (theRE UMTP_re) Utype() string {
	return "re"
}

func (theRE UMTP_re) Matchok(theStr string) bool {
	return (*regexp.Regexp)(&theRE).MatchString(theStr)
}

func matchtoken_prepare(sToken string) (UMTP, error) {
	if giMatchMode == MatchMode_RegExp {
		re, err := regexp.Compile(sToken)
		if err != nil {
			return nil, err
		}
		return UMTP_re(*re), nil
	}
	sP := match_prepare(sToken)
	sPR := UMTP_string(sP)
	return sPR, nil
}

func test_mtp() {
	const SearchToken = ".*st.*"
	const CheckString = "testme"
	mtp, err := matchtoken_prepare(SearchToken)
	if err != nil {
		return
	}
	fmt.Printf("%v:INFO:T MTP: %v ~ %v:%v\n", PRG_TAG, SearchToken, CheckString, mtp.Matchok(CheckString))
}

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
	if giMatchMode == MatchMode_RegExp {
		return match_regexp(sToCheck, sMatchTokenP)
	}
	return match_contains(sToCheck, sMatchTokenP)
}
