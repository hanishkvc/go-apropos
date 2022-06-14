// Maintain Maps wrt info about the Packages and their symbols
// HanishKVC, 2022

package main

import (
	"fmt"
)

type DBEntry struct {
	symbols map[string]string
	paths   []string
	cmts    []string
}

type TheDB map[string]DBEntry

var gDB TheDB

func identsmap_update(theMap map[string]string, identName string, identDoc string, identIsExported bool) {
	if identIsExported || gbAllSymbols {
		identDocCur, ok := theMap[identName]
		if !ok {
			theMap[identName] = identDoc
		} else {
			identDocCur = identDocCur + "; " + identDoc
			theMap[identName] = identDocCur
		}
	}
}

func db_add(theDB TheDB, pkgName string, path string, cmts string, idents map[string]string) {
	aPkg, ok := theDB[pkgName]
	if !ok {
		aPkg := DBEntry{}
		aPkg.symbols = idents
		aPkg.paths = make([]string, 0)
		aPkg.cmts = make([]string, 0)
		theDB[pkgName] = aPkg
	} else {
		for identName, identInfo := range idents {
			identsmap_update(aPkg.symbols, identName, identInfo, true)
		}
	}
	aPkg.paths = append(aPkg.paths, path)
	aPkg.cmts = append(aPkg.cmts, cmts)
}

func db_print() {
	for pkgName, identsMap := range gDBSymbols {
		fmt.Println(pkgName, identsMap)
	}
}

func db_print_pkgs() {
	for pkgName := range gDBSymbols {
		fmt.Printf("Package:%v:%v\n", pkgName, gDBPaths[pkgName])
	}
}

type MatchingPkgs map[string][]string

func matchingpkgs_add(thePkgs MatchingPkgs, pkgName string, id string) {
	_, ok := thePkgs[pkgName]
	if !ok {
		thePkgs[pkgName] = make([]string, 0)
	}
	thePkgs[pkgName] = append(thePkgs[pkgName], id)
	if giDEBUG > 10 {
		fmt.Printf("%v:DBUG:DB: MatchingPkgsAdd:%v:%v\n", PRG_TAG, pkgName, id)
	}
}

func db_find(sFind string, sFindCmt string) {
	if giDEBUG > 0 {
		fmt.Printf("\n%v:INFO: Possible matches for [%v] at [%v]\n", PRG_TAG, gFind, gBasePath)
	}
	pkgs := MatchingPkgs{}
	sFindP := match_prepare(sFind)
	sFindCmtP := match_prepare(sFindCmt)
	for pkgName, identsMap := range gDBSymbols {
		bFoundInPackage := false
		// Check symbols in the current package
		for id, idInfo := range identsMap {
			bFound := match_ok(id, sFindP) || match_ok(idInfo, sFindCmtP)
			if bFound {
				bFoundInPackage = true
				matchingpkgs_add(pkgs, pkgName, id)
			}
		}
		// If no match, check comments wrt current package
		if !bFoundInPackage && (gFindCmt != FINDCMT_DUMMY) {
			for _, cmt := range gDBCmts[pkgName] {
				if match_ok(cmt, sFindCmtP) {
					bFoundInPackage = true
				}
			}
			if bFoundInPackage {
				matchingpkgs_add(pkgs, pkgName, "???")
			}
		}
	}
	map_print(pkgs, " ", "\n")
}
