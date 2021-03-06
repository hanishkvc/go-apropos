// Maintain Maps wrt info about the Packages and their symbols
// HanishKVC, 2022

package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
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
			identData.Cmt = strings.TrimSpace(identData.Cmt)
			dbSymbols[identName] = identData
		} else {
			if identDataCur.Type != identData.Type {
				if identDataCur.Type == "" {
					identDataCur.Type = identData.Type
				} else {
					for _, c := range identData.Type {
						s := string(c)
						if !strings.Contains(identDataCur.Type, s) {
							identDataCur.Type += s
						}
					}
					if giDEBUG > 1 {
						fmt.Printf("%v:WARN:DBSymbolsUpdate: Mismatch in types [%v != %v] for %v\n", PRG_TAG, identDataCur.Type, identData.Type, identName)
					}
				}
			}
			sTrimmedCmt := strings.TrimSpace(identData.Cmt)
			if len(sTrimmedCmt) > 0 {
				identDataCur.Cmt += ("\n" + sTrimmedCmt)
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

func dbfilter_pkgs(theDB TheDB, matchPkgName string, matchMode MatchMode, caseSensitive bool) (TheDB, error) {
	newDB := TheDB{}
	pkgNameMatcher := New_Matcher(matchMode, matchPkgName, caseSensitive)
	for pkgName, pkgInfo := range theDB {
		if !pkgNameMatcher.Matchok(pkgName) {
			continue
		}
		db_add(newDB, pkgName, pkgInfo.Paths, pkgInfo.Cmts, pkgInfo.Symbols)
	}
	return newDB, nil
}

func _dbprint_pkgname(pkgNamePlus string) string {
	_, pkgName := filepath.Split(pkgNamePlus)
	sPkgName := fmt.Sprintf("%-16s:", pkgName)
	return sPkgName
}

func dbprint_all(theDB TheDB, sNamePrefix, sNameSuffix, sInfoPrefix, sInfoSuffix, sEnd string) {
	pkgNames := []string{}
	for pkgName := range theDB {
		pkgNames = append(pkgNames, pkgName)
	}
	sort.Strings(pkgNames)
	for _, pkgName := range pkgNames {
		spacedPkgName := _dbprint_pkgname(pkgName)
		pkgInfo := theDB[pkgName]
		//fmt.Println(pkgName, pkgData)
		fmt.Printf("%v%v%v", sNamePrefix, pkgName, sNameSuffix)
		dbprint_pkgpaths(pkgInfo.Paths, sInfoPrefix+spacedPkgName+"path:", sInfoSuffix)
		dbprint_pkgsymbols(pkgInfo.Symbols, sInfoPrefix+spacedPkgName, sInfoSuffix)
		fmt.Printf("%v", sEnd)
	}
}

func dbprint_pkgpaths(pkgPaths []string, sPrefix, sSuffix string) {
	paths := []string{}
	paths = append(paths, pkgPaths...)
	sort.Strings(paths)
	for _, path := range paths {
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

const gbForcedCommentsSplit = false

func _dprint_comments(sCmts string) {
	for _, curFS := range strings.Split(sCmts, "\n") {
		if !gbForcedCommentsSplit {
			fmt.Printf("\t%v\n", curFS)
			continue
		}
		for iS := 0; iS < len(curFS); iS += 70 {
			iE := iS + 70
			if iE > len(curFS) {
				iE = len(curFS)
			}
			fmt.Printf("\t%v\n", curFS[iS:iE])
		}
	}
	fmt.Println()
}

func dbprint_pkgsymbols(pkgSymbols DBSymbols, sPrefix, sSuffix string) {
	syms := []string{}
	for sym := range pkgSymbols {
		syms = append(syms, sym)
	}
	sort.Strings(syms)
	for _, sym := range syms {
		symInfo := pkgSymbols[sym]
		sCmt := symInfo.Cmt
		sCmt = strings.TrimSpace(sCmt)
		sCmt = strings.ReplaceAll(sCmt, "\n", " ")
		symPrint := fmt.Sprintf("%v:%-16s:%.80s", symInfo.Type, sym, sCmt)
		fmt.Printf("%v%v%v", sPrefix, symPrint, sSuffix)
		if gbFullComments {
			_dprint_comments(symInfo.Cmt)
		}
	}
}

func dbprint_symbols(theDB TheDB, sNamePrefix, sNameSuffix, sSymPrefix, sSymSuffix, sEnd string) {
	for pkgName, pkgInfo := range theDB {
		fmt.Printf("%v%v%v", sNamePrefix, pkgName, sNameSuffix)
		dbprint_pkgsymbols(pkgInfo.Symbols, sSymPrefix, sSymSuffix)
		fmt.Printf("%v", sEnd)
	}
}

func db_find(theDB TheDB, sFind string, sFindCmt string, sFindPkg string, matchMode MatchMode, caseSensitive bool, sortedResult bool) {
	if giDEBUG > 0 {
		fmt.Printf("\n%v:INFO: Possible matches for [%v] at [%v]\n", PRG_TAG, gFind, gBasePath)
	}
	matchingPkgs := make(TheDB)
	sFindP := New_Matcher(matchMode, sFind, caseSensitive)
	sFindCmtP := New_Matcher(matchMode, sFindCmt, caseSensitive)
	sFindPkgP := New_Matcher(matchMode, sFindPkg, caseSensitive)
	for pkgName, pkgData := range theDB {
		spacedPkgName := _dbprint_pkgname(pkgName)
		matchingSymbols := make(DBSymbols)
		bPkgNamePrinted := false
		// Honor any findpkg based package filtering
		if gFindPkg != FINDPKG_DEFAULT {
			if !sFindPkgP.Matchok(pkgName) {
				continue
			}
			if sortedResult {
				db_add(matchingPkgs, pkgName, pkgData.Paths, pkgData.Cmts, DBSymbols{})
			} else {
				bPkgNamePrinted = true
				fmt.Printf("\nPackage:%v\n", pkgName)
				dbprint_pkgpaths(theDB[pkgName].Paths, spacedPkgName+"path:", "\n")
			}
		}
		bFoundInPackage := false
		// Check symbols in the current package
		for id, idInfo := range pkgData.Symbols {
			bFound := sFindP.Matchok(id) || sFindCmtP.Matchok(idInfo.Cmt)
			if bFound {
				bFoundInPackage = true
				dbsymbols_update(matchingSymbols, id, idInfo, true)
			}
		}
		// If no match, check comments wrt current package
		if !bFoundInPackage && (gFindCmt != FINDCMT_DUMMY) {
			for _, cmt := range pkgData.Cmts {
				if sFindCmtP.Matchok(cmt) {
					bFoundInPackage = true
				}
			}
			if bFoundInPackage {
				dbsymbols_update(matchingSymbols, "???", DBSymbolInfo{"", "P"}, true)
			}
		}
		if bFoundInPackage {
			if sortedResult {
				db_add(matchingPkgs, pkgName, []string{}, []string{}, matchingSymbols)
			} else {
				if !bPkgNamePrinted {
					fmt.Printf("\nPackage:%v\n", pkgName)
				}
				dbprint_pkgsymbols(matchingSymbols, spacedPkgName, "\n")
			}
		}
	}
	if sortedResult {
		dbprint_all(matchingPkgs, "Package:", "\n", "", "\n", "\n")
	}
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
