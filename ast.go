// Parse the go source files to extract info about them
// HanishKVC, 2022

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

func test_dbsymbols_update() {
	aMap := map[string]DBSymbolInfo{"t100": {"doct of t100", "X"}}
	//bMap := DBSymbols{}
	dbsymbols_update(aMap, "t1", DBSymbolInfo{"doc for t1", ""}, true)
	fmt.Printf("%v:INFO:T MAPSTRUCT2: aMap:%v\n", PRG_TAG, aMap)
	dbsymbols_update(aMap, "t1", DBSymbolInfo{"doc for t1", ""}, true)
	dbsymbols_update(aMap, "t2", DBSymbolInfo{"doc for t2", ""}, true)
	fmt.Printf("%v:INFO:T MAPSTRUCT2: aMap after updates:%v\n", PRG_TAG, aMap)
}

type IdentyStats struct {
	identCnt uint
	funcCnt  uint
	valueCnt uint
	typeCnt  uint
}

var gIdentyStats IdentyStats

func (is IdentyStats) String() string {
	return fmt.Sprintf("{ i:%v f:%v v:%v t:%v }", is.identCnt, is.funcCnt, is.valueCnt, is.typeCnt)
}

func (is *IdentyStats) delta_summary() int {
	return (int(is.identCnt - (is.funcCnt + is.typeCnt + is.valueCnt)))
}

// Retreive info about the go source file specified
// It returns
//		the package name to which the file belongs
//		all the comments in the file
//		map of exported / all identifiers defined by the file
func gosrc_info(sFile string) (string, string, DBSymbols) {
	tfs := token.NewFileSet()
	astF, err := parser.ParseFile(tfs, sFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("%v:ERRR:AST: %v\n", PRG_TAG, err)
		return "", "", nil
	}
	pkgName := ""
	fileCmts := ""
	theIdents := DBSymbols{}
	genDeclCmt := ""
	genDeclType := "?"
	ast.Inspect(astF, func(n ast.Node) bool {
		bDigDeeper := true
		if n == nil {
			if giDEBUG > 6 {
				fmt.Printf("%v:INFO:AST: n:nil:: %v\n", PRG_TAG, n)
			}
			bDigDeeper = false // true or false may not matter here, but still
			return bDigDeeper
		}
		sType := "???"
		sExtra := ""
		switch t := n.(type) {
		/* Commenting out, bcas this path traps for both own as well as other packages' identifiers used by go source file being inspected
		case *ast.Ident: // This gives names of vars, consts, types and funcs also
			sType = "Identifier"
			sExtra = t.Name
			identsmap_update(theIdents, t.Name, "", t.IsExported())
			gIdentyStats.identCnt += 1
		*/
		case *ast.ValueSpec: // Through GenDecl
			sType = "ConstOrVar"
			sCmt := ":Cmt:" + t.Comment.Text() + ":Doc:" + t.Doc.Text()
			valCmt := genDeclCmt + "\n" + t.Comment.Text() + "\n" + t.Doc.Text()
			saExtra := []string{}
			for _, ident := range t.Names {
				saExtra = append(saExtra, ident.Name)
				dbsymbols_update(theIdents, ident.Name, DBSymbolInfo{valCmt, genDeclType}, ident.IsExported())
				gIdentyStats.valueCnt += 1
			}
			sExtra = "<" + strings.Join(saExtra, ",") + "> " + sCmt
		case *ast.TypeSpec: // Through GenDecl
			sType = "Type"
			sCmt := ":Cmt:" + t.Comment.Text() + ":Doc:" + t.Doc.Text()
			typeCmt := genDeclCmt + "\n" + t.Comment.Text() + "\n" + t.Doc.Text()
			sExtra = "<" + t.Name.Name + "> " + sCmt
			dbsymbols_update(theIdents, t.Name.Name, DBSymbolInfo{typeCmt, "T"}, t.Name.IsExported())
			gIdentyStats.typeCnt += 1
		case *ast.GenDecl:
			sType = "ImpTypeConstVar"
			genDeclType = "*"
			switch t.Tok {
			case token.IMPORT:
				sType = "GenDecl:Import"
				genDeclType = "I"
				bDigDeeper = false
			case token.CONST:
				sType = "GenDecl:Const"
				genDeclType = "C"
			case token.TYPE:
				sType = "GenDecl:Type"
				genDeclType = "T"
			case token.VAR:
				sType = "GenDecl:Var"
				genDeclType = "V"
			}
			/*
				for _, spec := range t.Specs {
					fmt.Printf("%v:DBUG:AST GenDecl:%v\n", PRG_TAG, spec.Doc.Text())
				}
			*/
			sExtra = t.Doc.Text()
			genDeclCmt = t.Doc.Text()
		case *ast.FuncDecl:
			sType = "Function"
			sExtra = t.Name.Name + ", :Doc:" + t.Doc.Text()
			dbsymbols_update(theIdents, t.Name.Name, DBSymbolInfo{t.Doc.Text(), "F"}, t.Name.IsExported())
			gIdentyStats.funcCnt += 1
		case *ast.Package: // Working on individual Go src files doesnt seem to encounter this type
			sType = "Package"
			sExtra = t.Name
			pkgName = t.Name
		case *ast.File: // This could give useful info
			sType = "File"
			sExtra = t.Name.Name
			pkgName = t.Name.Name
			//fmt.Printf("%v:INFO:AST: File.Scope: %v\n", PRG_TAG, t.Scope)
			for _, cmtG := range t.Comments {
				fileCmts += cmtG.Text()
			}
			sExtra += (", " + fileCmts)
		default:
			if giDEBUG > 6 {
				sType = "?" + reflect.TypeOf(t).String() + "?"
				//sExtra = t1.Name() + ";" + t1.Elem().Name()
				sExtra = reflect.ValueOf(t).String()
			}
		}
		if giDEBUG > 6 {
			fmt.Printf("%v:INFO:AST: n:%v:%v: %v\n", PRG_TAG, sType, sExtra, n)
		}
		return bDigDeeper
	})
	if giDEBUG > 5 {
		fmt.Printf("%v:INFO:AST: GoFile:%v:%v\n", PRG_TAG, pkgName, theIdents)
	}
	if giDEBUG > 20 {
		fmt.Printf("%v:DBUG:AST: GoFile:%v:%v:%v\n", PRG_TAG, sFile, gIdentyStats, gIdentyStats.delta_summary())
	}
	return pkgName, fileCmts, theIdents
}
