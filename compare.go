// logic to compare tokens
// HanishKVC, 2022

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const MATCHMODE_CONTAINS = "contains"
const MATCHMODE_REGEXP = "regexp"

type MatchMode uint

const (
	MatchMode_Contains MatchMode = iota
	MatchMode_RegExp
	MatchMode_Invalid
)

func matchmode_fromstr(sMode string) MatchMode {
	switch sMode {
	case MATCHMODE_CONTAINS:
		return MatchMode_Contains
	case MATCHMODE_REGEXP:
		return MatchMode_RegExp
	}
	fmt.Printf("%v:ERRR:MatchMode %v is unknown, exiting...\n", PRG_TAG, sMode)
	os.Exit(20)
	return MatchMode_Invalid // program wont reach here, just to keep go tools happy
}

func matchmode_tostr(mode MatchMode) string {
	switch mode {
	case MatchMode_Contains:
		return MATCHMODE_CONTAINS
	case MatchMode_RegExp:
		return MATCHMODE_REGEXP
	}
	fmt.Printf("%v:ERRR:MatchMode %v is unknown, exiting...\n", PRG_TAG, mode)
	os.Exit(20)
	return "ERROR:UNKNOWN" // program wont reach here, just to keep go tools happy
}

type Matcher_string string
type Matcher_re regexp.Regexp
type Matcher interface {
	Utype() string       // get the type of the matcher
	Matchok(string) bool // check if the given string matches the pattern registered with matcher
	Pattern() string     // retreive the string pattern registered with the matcher
}

func (o Matcher_string) Utype() string {
	return "string"
}

func (subStr Matcher_string) Matchok(theStr string) bool {
	//fmt.Printf("%v:INFO:MatcherString: is [%v] in [%v]\n", PRG_TAG, subStr, theStr)
	theStr = match_prepare(theStr)
	return strings.Contains(theStr, string(subStr))
}

func (subStr Matcher_string) Pattern() string {
	return string(subStr)
}

func (theRE *Matcher_re) Utype() string {
	return "re"
}

func (theRE *Matcher_re) Matchok(theStr string) bool {
	//fmt.Printf("%v:INFO:MatcherRE: is [%v] in [%v]\n", PRG_TAG, theRE, theStr)
	theStr = match_prepare(theStr)
	return (*regexp.Regexp)(theRE).MatchString(theStr)
}

func (theRE *Matcher_re) Pattern() string {
	return (*regexp.Regexp)(theRE).String()
}

// Based on current match mode either create
// 		a string based contains matcher
//		or a regexp based matcher
// The matcher takes care of case sensitivity wrt matching.
//		if case insensitive match is requested, currently it uses a simple to upper case conversion
//		irrespective of the type of matcher used
func matcher_create(pattern string) Matcher {
	if giMatchMode == MatchMode_RegExp {
		re := regexp.MustCompile(match_prepare(pattern))
		return Matcher_re(*re)
	}
	sP := match_prepare(pattern)
	sPR := Matcher_string(sP)
	return sPR
}

// Prepare a token / string for use by match_ok logic
func match_prepare(sToken string) string {
	if gbCaseSensitive {
		return sToken
	}
	return strings.ToUpper(sToken)
}
