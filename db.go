// Maintain Maps wrt info about the Packages and their symbols
// HanishKVC, 2022

package main

import (
	"fmt"
)

type DBSymbolInfo struct {
	Cmt  string `json:"c"`
	Type string `json:"t"`
}

type DBSymbols map[string]DBSymbolInfo

type DBEntry struct {
	Symbols DBSymbols `json:"s"`
	Paths   []string  `json:"p"`
	Cmts    []string  `json:"c"`
}

type TheDB map[string]DBEntry

var gDB TheDB = make(TheDB)

func dbsymbols_update(dbSymbols DBSymbols, identName string, identData DBSymbolInfo, identIsExported bool) {
	if identIsExported || gbAllSymbols {
		identDataCur, ok := dbSymbols[identName]
		if !ok {
			dbSymbols[identName] = identData
		} else {
			identDataCur.Cmt += ("; " + identData.Cmt)
			if identData.Type == "" {
				identDataCur.Cmt += ("; " + identData.Cmt)
			}
			dbSymbols[identName] = identDataCur
		}
	}
}

func db_add(theDB TheDB, pkgName string, pathS []string, cmtS []string, idents DBSymbols) {
	aPkg, ok := theDB[pkgName]
	if !ok {
		aPkg = DBEntry{}
		aPkg.Symbols = idents
		aPkg.Paths = make([]string, 0)
		aPkg.Cmts = make([]string, 0)
	} else {
		for identName, identInfo := range idents {
			dbsymbols_update(aPkg.Symbols, identName, identInfo, true)
		}
	}
	aPkg.Paths = append(aPkg.Paths, pathS...)
	aPkg.Cmts = append(aPkg.Cmts, cmtS...)
	theDB[pkgName] = aPkg
}

func dbfilter_pkgs(theDB TheDB, matchPkgName string) TheDB {
	newDB := TheDB{}
	for pkgName, pkgInfo := range theDB {
		if !match_ok(pkgName, matchPkgName) {
			continue
		}
		db_add(newDB, pkgName, pkgInfo.Paths, pkgInfo.Cmts, pkgInfo.Symbols)
	}
	return newDB
}

func dbprint_all(theDB TheDB, sNamePrefix, sNameSuffix, sInfoPrefix, sInfoSuffix, sEnd string) {
	for pkgName, pkgInfo := range theDB {
		//fmt.Println(pkgName, pkgData)
		fmt.Printf("%v%v%v", sNamePrefix, pkgName, sNameSuffix)
		dbprint_pkgpaths(pkgInfo.Paths, sInfoPrefix+"path:", sInfoSuffix)
		dbprint_pkgsymbols(pkgInfo.Symbols, sInfoPrefix+"sym:", sInfoSuffix)
		fmt.Printf("%v", sEnd)
	}
}

func dbprint_pkgpaths(pkgPaths []string, sPrefix, sSuffix string) {
	for _, path := range pkgPaths {
		fmt.Printf("%v%v%v", sPrefix, path, sSuffix)
	}
}

func dbprint_paths(theDB TheDB, sNamePrefix, sNameSuffix, sPathPrefix, sPathSuffix, sEnd string) {
	for pkgName, pkgInfo := range theDB {
		//fmt.Printf("%v%v%v%v%v", sNamePrefix, pkgName, sNameSuffix, theDB[pkgName].Paths, sEnd)
		fmt.Printf("%v%v%v", sNamePrefix, pkgName, sNameSuffix)
		dbprint_pkgpaths(pkgInfo.Paths, sPathPrefix, sPathSuffix)
		fmt.Printf("%v", sEnd)
	}
}

func dbprint_pkgsymbols(pkgSymbols DBSymbols, sPrefix, sSuffix string) {
	for sym, symInfo := range pkgSymbols {
		symPrint := symInfo.Type + ":" + sym
		fmt.Printf("%v%v%v", sPrefix, symPrint, sSuffix)
	}
}

func dbprint_symbols(theDB TheDB, sNamePrefix, sNameSuffix, sSymPrefix, sSymSuffix, sEnd string) {
	for pkgName, pkgInfo := range theDB {
		fmt.Printf("%v%v%v", sNamePrefix, pkgName, sNameSuffix)
		dbprint_pkgsymbols(pkgInfo.Symbols, sSymPrefix, sSymSuffix)
		fmt.Printf("%v", sEnd)
	}
}

func db_find(theDB TheDB, sFind string, sFindCmt string, sFindPkg string) {
	if giDEBUG > 0 {
		fmt.Printf("\n%v:INFO: Possible matches for [%v] at [%v]\n", PRG_TAG, gFind, gBasePath)
	}
	matchingPkgs := make(TheDB)
	sFindP := match_prepare(sFind)
	sFindCmtP := match_prepare(sFindCmt)
	sFindPkgP := match_prepare(sFindPkg)
	for pkgName, pkgData := range theDB {
		matchingSymbols := make(DBSymbols)
		// Honor any findpkg based package filtering
		if gFindPkg != FINDPKG_DEFAULT {
			if !match_ok(pkgName, sFindPkgP) {
				continue
			}
			if gbSortedResult {
				db_add(matchingPkgs, pkgName, pkgData.Paths, pkgData.Cmts, DBSymbols{})
			} else {
				fmt.Printf("Package:%v\n", pkgName)
				dbprint_pkgpaths(theDB[pkgName].Paths, "\tpath:", "\n")
			}
		}
		bFoundInPackage := false
		// Check symbols in the current package
		for id, idInfo := range pkgData.Symbols {
			bFound := match_ok(id, sFindP) || match_ok(idInfo.Cmt, sFindCmtP)
			if bFound {
				bFoundInPackage = true
				dbsymbols_update(matchingSymbols, id, idInfo, true)
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
				dbsymbols_update(matchingSymbols, "???", DBSymbolInfo{"", "P"}, true)
			}
		}
		if bFoundInPackage {
			if gbSortedResult {
				db_add(matchingPkgs, pkgName, []string{}, []string{}, matchingSymbols)
			} else {
				fmt.Printf("Package:%v\n", pkgName)
				dbprint_pkgsymbols(matchingSymbols, "\tsym:", "\n")
			}
		}
	}
	if gbSortedResult {
		if gFindPkg != FINDPKG_DEFAULT {
			dbprint_paths(matchingPkgs, "Package:", "\n", "\tpath:", "\n", "\n")
		}
		dbprint_symbols(matchingPkgs, "Package:", "\n", "\tsym:", "\n", "\n")
		//dbprint_symbols(matchingPkgs, "", " [", " ", " ", "]\n")
	}
	dbprint_all(matchingPkgs, "Package:", "\n", "\t", "\n", "\n")
}

func db_sane(theDB TheDB) bool {
	pkgCnt := 0
	symCnt := 0
	for _, pkgInfo := range theDB {
		pkgCnt += 1
		symCnt += len(pkgInfo.Symbols)
	}
	if giDEBUG > 2 {
		fmt.Printf("%v:INFO:DB: pkgCnt:%v, symCnt:%v\n", PRG_TAG, pkgCnt, symCnt)
	}
	if (pkgCnt < 1) || (symCnt < 1) {
		return false
	}

	return true
}
