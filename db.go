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

func db_find(sFind string) {
	pkgs := []string{}
	sFU := strings.ToUpper(sFind)
	for k, v := range gDB {
		for id, _ := range v {
			idU := strings.ToUpper(id)
			if strings.Contains(idU, sFU) {
				pkgs = append(pkgs, k)
			}
		}
	}
	for _, pkg := range pkgs {
		fmt.Println(pkg)
	}
}
