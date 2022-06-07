// Parse the go source files to extract info about them
// HanishKVC, 2022

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func gosrc_info(sFile string) (string, map[string]int) {
	tfs := token.NewFileSet()
	astF, err := parser.ParseFile(tfs, sFile, nil, parser.ParseComments)
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
			sType = "Comment"
			sExtra = t.Text
		case *ast.CommentGroup: // Dont seem to encounter this
			sType = "CommentGroup"
			sExtra = t.Text()
		case *ast.Ident: // This gives names of vars, consts and funcs also
			sType = "Identifier"
			sExtra = t.Name
			if t.IsExported() || gbALL {
				theIdents[t.Name] += 1
			}
		case *ast.GenDecl:
			sType = "ImpTypeConstVar"
			sExtra = t.Doc.Text()
			if (len(sExtra) > 0) && (giDEBUG > 10) {
				fmt.Printf("%v:DBUG:AST GenDecl:cmt: %v:%v\n", PRG_TAG, t.Specs, sExtra)
			}
		case *ast.FuncDecl:
			sType = "Function"
			sExtra = t.Name.Name
			sComment := t.Doc.Text()
			if (len(sComment) > 0) && (giDEBUG > 10) {
				fmt.Printf("%v:DBUG:AST FuncDecl:cmt: %v:%v\n", PRG_TAG, sExtra, sComment)
			}
		case *ast.Package: // Dont seem to encounter this type
			sType = "Package"
			sExtra = t.Name
			pkgName = t.Name
		case *ast.File: // This could give useful info
			sType = "File"
			sExtra = t.Name.Name
			pkgName = t.Name.Name
			//fmt.Printf("%v:INFO:AST: File.Scope: %v\n", PRG_TAG, t.Scope)
			if giDEBUG > 10 {
				for _, cmtG := range t.Comments {
					fmt.Printf("%v:DBUG:AST File:cmt: %v\n", PRG_TAG, cmtG.Text())
				}
			}
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
