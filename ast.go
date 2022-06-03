package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func gosrc_info(sFile string) (string, map[string]int) {
	tfs := token.NewFileSet()
	astF, err := parser.ParseFile(tfs, sFile, nil, 0)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", nil
	}
	pkgName := ""
	theIdents := map[string]int{}
	ast.Inspect(astF, func(n ast.Node) bool {
		sType := "???"
		sExtra := ""
		switch t := n.(type) {
		case *ast.CommentGroup: // Dont seem to encounter this
			sExtra = t.Text()
		case *ast.Ident: // This gives names of vars, consts and funcs also
			sType = "Identifier"
			sExtra = t.Name
			theIdents[t.Name] += 1
		case *ast.FuncDecl:
			sType = "Function"
			sExtra = t.Name.Name
		case *ast.Package: // Dont seem to encounter this type
			sType = "Package"
			sExtra = t.Name
			pkgName = t.Name
		case *ast.File: // This could give useful info
			sType = "File"
			sExtra = t.Name.Name
			pkgName = t.Name.Name
		default:
			//t1 := reflect.TypeOf(t)
			//sExtra = t1.Name()
		}
		if giDEBUG > 6 {
			fmt.Printf("n:%v:%v: %v\n", sType, sExtra, n)
		}
		return true
	})
	return pkgName, theIdents
}
