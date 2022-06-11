// Help build / cache the db of identifiers++
// HanishKVC, 2022

package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

var gCacheFile = "~/.cache/goapropos.db"

func handle_file(sFile string) {
	if !strings.HasSuffix(sFile, "go") {
		return
	}
	if strings.HasSuffix(sFile, "_test.go") {
		return
	}
	for _, mpath := range gSkipFiles {
		if strings.Contains(sFile, mpath) {
			return
		}
	}
	name, cmts, idents := gosrc_info(sFile)
	if gFindPkg != FINDPKG_DEFAULT {
		if !match_ok(name, gFindPkgP) {
			return
		}
	}
	db_add(name, sFile, cmts, idents)
}

func do_walkdir(sPath string) {
	oFS := os.DirFS(sPath)
	if giDEBUG > 10 {
		fmt.Printf("%v:INFO:WALKDIR: oFS: %v\n", PRG_TAG, oFS)
	}
	fs.WalkDir(oFS, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%v:ERRR:WALKDIR: path: %v, Err: %v\n", PRG_TAG, path, err)
			return err
		}
		var sPType string
		deT := de.Type()
		if deT.IsDir() {
			sPType = "Dir"
		} else if deT.IsRegular() {
			sPType = "File"
		} else {
			sPType = "???"
		}
		if giDEBUG > 2 {
			fmt.Printf("%v:INFO:WALKDIR: %v:path: %v\n", PRG_TAG, sPType, path)
		}
		if sPType == "File" {
			theFile := sPath + string(os.PathSeparator) + path
			handle_file(theFile)
		}
		return nil
	})
}

func prep_dir(sPath string) {
}

func adjust_path(sPath string) (string, error) {
	if !strings.HasPrefix(sPath, "~/") {
		return sPath, nil
	}
	sHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: HomeDir:%v\n", PRG_TAG, err)
		return "", err
	}
	sPath = strings.Replace(sPath, "~", sHomeDir, 1)
	return sPath, nil
}

func save_dbs() error {
	sDB, err := json.Marshal(gDB)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:Marshal:%v\n", PRG_TAG, err)
		return err
	}
	sCacheFile, err := adjust_path(gCacheFile)
	if err != nil {
		return err
	}
	err = os.WriteFile(sCacheFile, sDB, 0400)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:WriteFile:%v\n", PRG_TAG, err)
		return err
	}
	fmt.Println("DBUG:JSON", string(sDB))
	return nil
}
