package main

import "fmt"

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
