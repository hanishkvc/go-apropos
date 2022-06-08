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

func identsmap_update(theMap map[string]Ident, identName, identDoc string, identIsExported bool) {
	if identIsExported || gbALL {
		ident, ok := theMap[identName]
		if !ok {
			theMap[identName] = Ident{1, identDoc}
		} else {
			ident.cnt += 1
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
			ident, ok := gDB[name][k]
			if !ok {
				gDB[name][k] = Ident{v.cnt, v.doc}
			} else {
				ident.cnt += v.cnt
				ident.doc += (";" + v.doc)
				gDB[name][k] = ident // Do I need this? Need to check ie is ident a reference or a copy
			}
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
	for pkg, ids := range pkgs {
		fmt.Println(pkg, ids)
	}
}
