package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func gosrc_info(sFile string) map[string]string {
	tfs := token.NewFileSet()
	astF, err := parser.ParseFile(tfs, sFile, nil, 0)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	theMap := map[string]string{}
	ast.Inspect(astF, func(n ast.Node) bool {
		fmt.Printf("n: %v\n", n)
		return true
	})
	return theMap
}
