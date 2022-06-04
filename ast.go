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
		fmt.Printf("%v:ERRR:AST: %v\n", PRG_TAG, err)
		return "", nil
	}
	pkgName := ""
	theIdents := map[string]int{}
	ast.Inspect(astF, func(n ast.Node) bool {
		sType := "???"
		sExtra := ""
		switch t := n.(type) {
		case *ast.Comment:
			fmt.Printf("%v:INFO:AST: Comment:%v\n", PRG_TAG, t.Text)
		case *ast.CommentGroup: // Dont seem to encounter this
			sExtra = t.Text()
		case *ast.Ident: // This gives names of vars, consts and funcs also
			sType = "Identifier"
			sExtra = t.Name
			if t.IsExported() || gbALL {
				theIdents[t.Name] += 1
			}
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
			//fmt.Printf("%v:INFO:AST: File.Scope: %v\n", PRG_TAG, t.Scope)
		default:
			//t1 := reflect.TypeOf(t)
			//sExtra = t1.Name()
		}
		if giDEBUG > 6 {
			fmt.Printf("%v:INFO:AST: n:%v:%v: %v\n", PRG_TAG, sType, sExtra, n)
		}
		return true
	})
	if giDEBUG > 5 {
		fmt.Printf("%v:INFO:AST: GoFile:%v:%v\n", PRG_TAG, pkgName, theIdents)
	}
	return pkgName, theIdents
}
