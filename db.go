package main

import (
	"fmt"
	"strings"
)

var gDB = make(map[string]map[string]int)

func db_add(name string, idents map[string]int) {
	_, ok := gDB[name]
	if !ok {
		gDB[name] = idents
		return
	}
	for k, v := range idents {
		gDB[name][k] += v
	}
}

func db_print() {
	for k, v := range gDB {
		fmt.Println(k, v)
	}
}

func db_print_pkgs() {
	for k := range gDB {
		fmt.Printf("Package: %v\n", k)
	}
}

func db_find(sFind string) {
	if giDEBUG > 0 {
		fmt.Printf("\n%v:INFO: Possible matches for [%v] at [%v]\n", PRG_TAG, gFind, gBasePath)
	}
	pkgs := map[string][]string{}
	sFU := strings.ToUpper(sFind)
	for k, v := range gDB {
		for id, _ := range v {
			idU := strings.ToUpper(id)
			if strings.Contains(idU, sFU) {
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
