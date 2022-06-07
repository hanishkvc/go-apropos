// Parse the go source files to extract info about them
// HanishKVC, 2022

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
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
		case *ast.Ident: // This gives names of vars, consts and funcs also
			sType = "Identifier"
			sExtra = t.Name
			if t.IsExported() || gbALL {
				theIdents[t.Name] += 1
			}
		case *ast.ValueSpec:
			sType = "ConstOrVar"
			saExtra := []string{}
			for _, ident := range t.Names {
				saExtra = append(saExtra, ident.Name)
				//fmt.Printf("%v:DBUG:AST ValueSpec:%v:%v\n", PRG_TAG, ident.Name, ident.String())
			}
			sExtra = "<" + strings.Join(saExtra, ",") + "> :Cmt:" + t.Comment.Text() + "__AND__" + t.Doc.Text()
		case *ast.TypeSpec:
			sType = "Type"
			sExtra = "<" + t.Name.Name + "> :Cmt:" + t.Comment.Text() + "__AND__" + t.Doc.Text()
		case *ast.GenDecl:
			sType = "ImpTypeConstVar"
			switch t.Tok {
			case token.IMPORT:
				sType = "GenDecl:Import"
			case token.CONST:
				sType = "GenDecl:Const"
			case token.TYPE:
				sType = "GenDecl:Type"
			case token.VAR:
				sType = "GenDecl:Var"
			}
			/*
				for _, spec := range t.Specs {
					fmt.Printf("%v:DBUG:AST GenDecl:%v\n", PRG_TAG, spec.Doc.Text())
				}
			*/
			sExtra = t.Doc.Text()
		case *ast.FuncDecl:
			sType = "Function"
			sExtra = t.Name.Name + ", :Cmt:" + t.Doc.Text()
		case *ast.Package: // Dont seem to encounter this type
			sType = "Package"
			sExtra = t.Name
			pkgName = t.Name
		case *ast.File: // This could give useful info
			sType = "File"
			sExtra = t.Name.Name
			pkgName = t.Name.Name
			//fmt.Printf("%v:INFO:AST: File.Scope: %v\n", PRG_TAG, t.Scope)
			if giDEBUG > 20 {
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
