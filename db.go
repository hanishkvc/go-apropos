// Maintain Maps wrt info about the Packages and their symbols
// HanishKVC, 2022

package main

import (
	"encoding/json"
	"fmt"
)

type Ident struct {
	cnt int
	doc string
	// name string
}

var gDB = make(map[string]map[string]Ident)
var gDBPaths = make(map[string][]string)
var gDBCmts = make(map[string][]string)

func (o Ident) MarshalJSON() ([]byte, error) {
	docJSONB, err := json.Marshal(o.doc)
	if err != nil {
		fmt.Printf("%v:ERRR:DB: IdentJSON:%v\n", PRG_TAG, err)
		return nil, err
	}
	docJSON := string(docJSONB)
	fmt.Printf("docJSONB: %v\n", docJSONB)
	fmt.Printf("docJSON: %v\n", docJSON)
	identJSON := fmt.Sprintf("{ %v: %v }", docJSON, o.cnt)
	identJSONB := []byte(identJSON)
	fmt.Printf("IdentJSon: %v\n", identJSON)
	fmt.Printf("identJSONB: %v\n", identJSONB)
	return identJSONB, nil
}

func identsmap_update(theMap map[string]Ident, identName string, identCnt int, identDoc string, identIsExported bool) {
	if identIsExported || gbAllSymbols {
		ident, ok := theMap[identName]
		if !ok {
			theMap[identName] = Ident{identCnt, identDoc}
		} else {
			ident.cnt += identCnt
			ident.doc = ident.doc + "; " + identDoc
			theMap[identName] = ident
		}
	}
}

func db_add(pkgName string, path string, cmts string, idents map[string]Ident) {
	_, ok := gDB[pkgName]
	if !ok {
		gDB[pkgName] = idents
		gDBPaths[pkgName] = make([]string, 0)
		gDBCmts[pkgName] = make([]string, 0)
	} else {
		for identName, identInfo := range idents {
			identsmap_update(gDB[pkgName], identName, identInfo.cnt, identInfo.doc, true)
		}
	}
	gDBPaths[pkgName] = append(gDBPaths[pkgName], path)
	gDBCmts[pkgName] = append(gDBCmts[pkgName], cmts)
}

func db_print() {
	for pkgName, identsMap := range gDB {
		fmt.Println(pkgName, identsMap)
	}
}

func db_print_pkgs() {
	for pkgName := range gDB {
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
	for pkgName, identsMap := range gDB {
		bFoundInPackage := false
		// Check symbols in the current package
		for id, idInfo := range identsMap {
			bFound := match_ok(id, sFindP) || match_ok(idInfo.doc, sFindCmtP)
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
