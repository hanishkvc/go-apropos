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

func identsmap_update(theMap map[string]Ident, identName string, identCnt int, identDoc string, identIsExported bool) {
	if identIsExported || gbALL {
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

func db_add(name string, path string, idents map[string]Ident) {
	_, ok := gDB[name]
	if !ok {
		gDB[name] = idents
		gDBPaths[name] = make([]string, 0)
	} else {
		for k, v := range idents {
			identsmap_update(gDB[name], k, v.cnt, v.doc, true)
		}
	}
	gDBPaths[name] = append(gDBPaths[name], path)
}

func db_print() {
	for k, v := range gDB {
		fmt.Println(k, v)
	}
}

func db_print_pkgs() {
	for k := range gDB {
		fmt.Printf("Package: %v, %v\n", k, gDBPaths[k])
	}
}

func db_find(sFind string) {
	if giDEBUG > 0 {
		fmt.Printf("\n%v:INFO: Possible matches for [%v] at [%v]\n", PRG_TAG, gFind, gBasePath)
	}
	pkgs := map[string][]string{}
	sFindP := match_prepare(sFind)
	for k, v := range gDB {
		for id, _ := range v {
			if match_ok(id, sFindP) {
				_, ok := pkgs[k]
				if !ok {
					pkgs[k] = make([]string, 0)
				}
				pkgs[k] = append(pkgs[k], id)
			}
		}
	}
	map_print(pkgs, " ", "\n")
}
