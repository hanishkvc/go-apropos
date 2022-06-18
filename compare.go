// logic to compare tokens
// HanishKVC, 2022

package main

import (
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

func matchmode_tostr(mode MatchMode) string {
	switch mode {
	case MatchMode_Contains:
		return MATCHMODE_CONTAINS
	case MatchMode_RegExp:
		return MATCHMODE_REGEXP
	}
	return "ERROR:UNKNOWN"
}

type Matcher_string string
type Matcher_re regexp.Regexp
type Matcher interface {
	Utype() string
	Matchok(string) bool
}

func (o Matcher_string) Utype() string {
	return "string"
}

func (subStr Matcher_string) Matchok(theStr string) bool {
	return strings.Contains(theStr, string(subStr))
}

func (theRE Matcher_re) Utype() string {
	return "re"
}

func (theRE Matcher_re) Matchok(theStr string) bool {
	return (*regexp.Regexp)(&theRE).MatchString(theStr)
}

func matcher_create(sToken string) (Matcher, error) {
	if giMatchMode == MatchMode_RegExp {
		re, err := regexp.Compile(match_prepare(sToken))
		if err != nil {
			return nil, err
		}
		return Matcher_re(*re), nil
	}
	sP := match_prepare(sToken)
	sPR := Matcher_string(sP)
	return sPR, nil
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
