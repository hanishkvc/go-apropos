// Help build / cache the db of identifiers++
// HanishKVC, 2022

package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"syscall"
)

var gCacheBase = "~/.cache"

const gDBCacheFile = "goapropos.db"
const gDBPathsCacheFile = "goapropos.dbpaths"
const gDBCmtsCacheFile = "goapropos.dbcmts"

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

func cache_filenames(cacheFile string) (string, error) {
	sCacheBase, err := adjust_path(gCacheBase)
	if err != nil {
		return "", err
	}
	sDBCacheFile := sCacheBase + string(os.PathSeparator) + cacheFile
	return sDBCacheFile, nil
}

func save_db(theDB any, cacheFile string) error {
	sDB, err := json.Marshal(theDB)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:Marshal:%v\n", PRG_TAG, err)
		return err
	}
	sDBCacheFile, err := cache_filenames(cacheFile)
	if err != nil {
		return err
	}
	err = os.WriteFile(sDBCacheFile, sDB, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:WriteFile:%v\n", PRG_TAG, err)
		return err
	}
	if giDEBUG > 20 {
		fmt.Printf("%v:DBUG:DB+: SaveDBs:gDB:JSON:%v\n", PRG_TAG, string(sDB))
	}
	return nil
}

func save_dbs() error {
	err := save_db(gDBSymbols, gDBCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:DBSymbols:%v\n", PRG_TAG, err)
		return err
	}
	err = save_db(gDBPaths, gDBPathsCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:DBPaths:%v\n", PRG_TAG, err)
		return err
	}
	err = save_db(gDBCmts, gDBCmtsCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:DBCmts:%v\n", PRG_TAG, err)
		return err
	}
	return nil
}

func load_db(theDB any, cacheFile string) error {
	sDBCacheFile, err := cache_filenames(cacheFile)
	if err != nil {
		return err
	}
	bsDB, err := os.ReadFile(sDBCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: LoadDBs:ReadFile:%v\n", PRG_TAG, err)
		return err
	}
	if giDEBUG > 20 {
		fmt.Printf("%v:DBUG:DB+: LoadDBs:Read: %v\n", PRG_TAG, string(bsDB))
		fmt.Printf("%v:DBUG:DB+: LoadDBs:gDB:Before: %v\n", PRG_TAG, theDB)
	}
	err = json.Unmarshal(bsDB, theDB)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: LoadDBs:Unmarshal:%v\n", PRG_TAG, err)
		return err
	}
	if giDEBUG > 20 {
		fmt.Printf("%v:DBUG:DB+: LoadDBs:gDB: %v\n", PRG_TAG, theDB)
	}
	return nil
}

func load_dbs() error {
	err := load_db(&gDBSymbols, gDBCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: LoadDBs:DBSymbols:%v\n", PRG_TAG, err)
		return err
	}
	err = load_db(&gDBPaths, gDBPathsCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: LoadDBs:DBPaths:%v\n", PRG_TAG, err)
		return err
	}
	err = load_db(&gDBCmts, gDBCmtsCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: LoadDBs:DBCmts:%v\n", PRG_TAG, err)
		return err
	}
	return nil
}
