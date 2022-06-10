// Maintain Maps wrt info about the Packages and their symbols
// HanishKVC, 2022

package main

import (
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

func db_find(sFind string, sFindCmt string) {
	if giDEBUG > 0 {
		fmt.Printf("\n%v:INFO: Possible matches for [%v] at [%v]\n", PRG_TAG, gFind, gBasePath)
	}
	pkgs := map[string][]string{}
	sFindP := match_prepare(sFind)
	sFindCmtP := match_prepare(sFindCmt)
	for pkgName, identsMap := range gDB {
		for id, idInfo := range identsMap {
			bFound := match_ok(id, sFindP) || match_ok(idInfo.doc, sFindCmtP)
			if bFound {
				_, ok := pkgs[pkgName]
				if !ok {
					pkgs[pkgName] = make([]string, 0)
				}
				pkgs[pkgName] = append(pkgs[pkgName], id)
			}
		}
	}
	map_print(pkgs, " ", "\n")
}
