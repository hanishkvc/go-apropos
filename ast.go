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

func test_identsmap_update() {
	aMap := map[string]Ident{}
	identsmap_update(aMap, "t1", 1, "doc for t1", true)
	fmt.Printf("%v:INFO:T MAPSTRUCT2: aMap:%v\n", PRG_TAG, aMap)
	identsmap_update(aMap, "t1", 1, "doc for t1", true)
	identsmap_update(aMap, "t2", 1, "doc for t2", true)
	fmt.Printf("%v:INFO:T MAPSTRUCT2: aMap after updates:%v\n", PRG_TAG, aMap)
}

func gosrc_info(sFile string) (string, map[string]Ident) {
	tfs := token.NewFileSet()
	astF, err := parser.ParseFile(tfs, sFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("%v:ERRR:AST: %v\n", PRG_TAG, err)
		return "", nil
	}
	pkgName := ""
	theIdents := map[string]Ident{}
	ast.Inspect(astF, func(n ast.Node) bool {
		bDigDeeper := true
		if n == nil {
			return bDigDeeper // true or false maynot matter here
		}
		sType := "???"
		sExtra := ""
		switch t := n.(type) {
		case *ast.Ident: // This gives names of vars, consts and funcs also
			sType = "Identifier"
			sExtra = t.Name
			identsmap_update(theIdents, t.Name, 1, "", t.IsExported())
		case *ast.ValueSpec:
			sType = "ConstOrVar"
			sCmt := ":Cmt:" + t.Comment.Text() + ":Doc:" + t.Doc.Text()
			saExtra := []string{}
			for _, ident := range t.Names {
				saExtra = append(saExtra, ident.Name)
				identsmap_update(theIdents, ident.Name, 1, sCmt, ident.IsExported())
			}
			sExtra = "<" + strings.Join(saExtra, ",") + "> " + sCmt
		case *ast.TypeSpec:
			sType = "Type"
			sCmt := ":Cmt:" + t.Comment.Text() + ":Doc:" + t.Doc.Text()
			sExtra = "<" + t.Name.Name + "> " + sCmt
			identsmap_update(theIdents, t.Name.Name, 1, sCmt, t.Name.IsExported())
		case *ast.GenDecl:
			sType = "ImpTypeConstVar"
			switch t.Tok {
			case token.IMPORT:
				sType = "GenDecl:Import"
				bDigDeeper = false
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
			sExtra = t.Name.Name + ", :Doc:" + t.Doc.Text()
			identsmap_update(theIdents, t.Name.Name, 1, t.Doc.Text(), t.Name.IsExported())
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
		return bDigDeeper
	})
	if giDEBUG > 5 {
		fmt.Printf("%v:INFO:AST: GoFile:%v:%v\n", PRG_TAG, pkgName, theIdents)
	}
	return pkgName, theIdents
}
