// Maintain Maps wrt info about the Packages and their symbols
// HanishKVC, 2022

package main

import (
	"fmt"
)

type SymbolEntry struct {
	Cmt  string `json:"c"`
	Type string `json:"t"`
}

type DBEntry struct {
	Symbols map[string]SymbolEntry `json:"s"`
	Paths   []string               `json:"p"`
	Cmts    []string               `json:"c"`
}

type TheDB map[string]DBEntry

var gDB TheDB = make(TheDB)

func identsmap_update(theMap map[string]SymbolEntry, identName string, identData SymbolEntry, identIsExported bool) {
	if identIsExported || gbAllSymbols {
		identDataCur, ok := theMap[identName]
		if !ok {
			theMap[identName] = identData
		} else {
			identDataCur.Cmt += ("; " + identData.Cmt)
			if identData.Type == "" {
				identDataCur.Cmt += ("; " + identData.Cmt)
			}
			theMap[identName] = identDataCur
		}
	}
}

func db_add(theDB TheDB, pkgName string, pathS []string, cmtS []string, idents map[string]SymbolEntry) {
	aPkg, ok := theDB[pkgName]
	if !ok {
		aPkg = DBEntry{}
		aPkg.Symbols = idents
		aPkg.Paths = make([]string, 0)
		aPkg.Cmts = make([]string, 0)
	} else {
		for identName, identInfo := range idents {
			identsmap_update(aPkg.Symbols, identName, identInfo, true)
		}
	}
	aPkg.Paths = append(aPkg.Paths, pathS...)
	aPkg.Cmts = append(aPkg.Cmts, cmtS...)
	theDB[pkgName] = aPkg
}

func dbprint_all_all(theDB TheDB) {
	for pkgName, pkgData := range theDB {
		fmt.Println(pkgName, pkgData)
	}
}

func dbprint_all_paths(theDB TheDB, bAllPkgs bool) {
	for pkgName := range theDB {
		if !bAllPkgs && (gFindPkg != FINDPKG_DEFAULT) {
			if !match_ok(pkgName, gFindPkgP) {
				continue
			}
		}
		fmt.Printf("Package:%v:%v\n", pkgName, theDB[pkgName].Paths)
	}
}

type MatchingPkgs map[string][]string

func matchingpkgs_add(thePkgs MatchingPkgs, pkgName string, datas []string) {
	_, ok := thePkgs[pkgName]
	if !ok {
		thePkgs[pkgName] = make([]string, 0)
	}
	thePkgs[pkgName] = append(thePkgs[pkgName], datas...)
	if giDEBUG > 10 {
		fmt.Printf("%v:DBUG:DB: MatchingPkgsAdd:%v:%v\n", PRG_TAG, pkgName, datas)
	}
}

func db_find(theDB TheDB, sFind string, sFindCmt string, sFindPkg string) {
	if giDEBUG > 0 {
		fmt.Printf("\n%v:INFO: Possible matches for [%v] at [%v]\n", PRG_TAG, gFind, gBasePath)
	}
	matchingPkgSymbols := MatchingPkgs{}
	matchingPkgPaths := MatchingPkgs{}
	sFindP := match_prepare(sFind)
	sFindCmtP := match_prepare(sFindCmt)
	sFindPkgP := match_prepare(sFindPkg)
	for pkgName, pkgData := range theDB {
		// Honor any findpkg based package filtering
		if gFindPkg != FINDPKG_DEFAULT {
			if !match_ok(pkgName, sFindPkgP) {
				continue
			}
			if gbSortedResult {
				matchingpkgs_add(matchingPkgPaths, pkgName, theDB[pkgName].Paths)
			} else {
				fmt.Printf("Package:%v:%v\n", pkgName, theDB[pkgName].Paths)
			}
		}
		bFoundInPackage := false
		// Check symbols in the current package
		for id, idInfo := range pkgData.Symbols {
			bFound := match_ok(id, sFindP) || match_ok(idInfo.Cmt, sFindCmtP)
			if bFound {
				bFoundInPackage = true
				idType := idInfo.Type + ":"
				matchingpkgs_add(matchingPkgSymbols, pkgName, []string{idType + id})
			}
		}
		// If no match, check comments wrt current package
		if !bFoundInPackage && (gFindCmt != FINDCMT_DUMMY) {
			for _, cmt := range pkgData.Cmts {
				if match_ok(cmt, sFindCmtP) {
					bFoundInPackage = true
				}
			}
			if bFoundInPackage {
				matchingpkgs_add(matchingPkgSymbols, pkgName, []string{"p:???"})
			}
		}
		if bFoundInPackage && !gbSortedResult {
			fmt.Printf("%v %v\n", pkgName, matchingPkgSymbols[pkgName])
		}
	}
	if gbSortedResult {
		map_print(matchingPkgPaths, "Package", ":", ":", "\n")
		map_print(matchingPkgSymbols, "", "", " ", "\n")
	}
}

func db_sane(theDB TheDB) bool {
	if len(theDB) <= 1 {
		return false
	}
	return true
}
