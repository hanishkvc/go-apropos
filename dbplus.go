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
	"time"
)

func handle_file(theDB TheDB, sFile string) {
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
	db_add(theDB, name, []string{sFile}, []string{cmts}, idents)
}

const GR_NOMORE = "__NO_MORE__"
const GR_COUNT = 16
const GR_CHANDEPTH = 256

type TrackHF struct {
	todo   int
	done   int
	hfChan chan string
	bOver  bool
	theDB  TheDB
}

var gTrackHFs [GR_COUNT]TrackHF

func gr_hf_start() {
	for i := 0; i < GR_COUNT; i++ {
		gTrackHFs[i].hfChan = make(chan string, GR_CHANDEPTH)
		gTrackHFs[i].theDB = make(TheDB, 0)
		go gr_handlefile(i)
	}
}

func gr_hf_stop() {
	// Inform all GRs about no more jobs
	for i := 0; i < GR_COUNT; i++ {
		gTrackHFs[i].hfChan <- GR_NOMORE
	}
	// Wait for all GRs to finish
	for {
		bAllDone := true
		for i := 0; i < GR_COUNT; i++ {
			if gTrackHFs[i].bOver {
				continue
			}
			bAllDone = false
			if gTrackHFs[i].todo == gTrackHFs[i].done {
				gTrackHFs[i].bOver = true
			}
			if giDEBUG > 1 {
				fmt.Printf("%v:INFO:GRHF:%v: Walk over, waiting for data:%v\n", PRG_TAG, i, gTrackHFs[i])
			}
		}
		if bAllDone {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	// Create the Merged DB
	for i := 0; i < GR_COUNT; i++ {
		for pkgName, pkgData := range gTrackHFs[i].theDB {
			db_add(gDB, pkgName, pkgData.Paths, pkgData.Cmts, pkgData.Symbols)
		}
	}
}

func gr_handlefile(i int) {
	bNoMore := false
	for !bNoMore {
		sFile := <-gTrackHFs[i].hfChan
		if sFile == GR_NOMORE {
			if giDEBUG > 1 {
				fmt.Printf("%v:INFO:GRHF:%v: NoMoreFiles\n", PRG_TAG, i)
			}
			break
		}
		handle_file(gTrackHFs[i].theDB, sFile)
		gTrackHFs[i].done += 1
	}
}

func do_walkdir(sPath string) {
	gr_hf_start()
	defer gr_hf_stop()
	fileCnt := 0
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
			iGR := fileCnt % GR_COUNT
			gTrackHFs[iGR].todo += 1
			gTrackHFs[iGR].hfChan <- theFile
			fileCnt += 1
		}
		return nil
	})
}

func prep_dir() error {
	sCacheBase, err := adjust_path(gCacheBase)
	if err != nil {
		return err
	}
	err = os.MkdirAll(sCacheBase, 0700)
	if err != nil {
		fmt.Printf("%v:ERRR:PrepDir: %v\n", PRG_TAG, err)
		return err
	}
	if giDEBUG > 2 {
		fmt.Printf("%v:INFO:PrepDir: %v, exists\n", PRG_TAG, sCacheBase)
	}
	return nil
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
	var sDB []byte
	var err error
	if gbIndentJSON {
		sDB, err = json.MarshalIndent(theDB, "", " ")
	} else {
		sDB, err = json.Marshal(theDB)
	}
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDB:Marshal:%v\n", PRG_TAG, err)
		return err
	}
	sDBCacheFile, err := cache_filenames(cacheFile)
	if err != nil {
		return err
	}
	err = os.WriteFile(sDBCacheFile, sDB, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDB:WriteFile:%v\n", PRG_TAG, err)
		return err
	}
	if giDEBUG > 20 {
		fmt.Printf("%v:DBUG:DB+: SaveDB:DB:JSON:%v\n", PRG_TAG, string(sDB))
	}
	return nil
}

func save_dbs() error {
	err := save_db(gDB, gDBAllCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: SaveDBs:gDB:%v\n", PRG_TAG, err)
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
		fmt.Printf("%v:ERRR:DB+: LoadDB:ReadFile:%v\n", PRG_TAG, err)
		return err
	}
	if giDEBUG > 20 {
		fmt.Printf("%v:DBUG:DB+: LoadDB:Read: %v\n", PRG_TAG, string(bsDB))
		fmt.Printf("%v:DBUG:DB+: LoadDB:DB:Before: %v\n", PRG_TAG, theDB)
	}
	err = json.Unmarshal(bsDB, theDB)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: LoadDB:Unmarshal:%v\n", PRG_TAG, err)
		return err
	}
	if giDEBUG > 20 {
		fmt.Printf("%v:DBUG:DB+: LoadDB:DB:After: %v\n", PRG_TAG, theDB)
	}
	return nil
}

func load_dbs() error {
	err := load_db(&gDB, gDBAllCacheFile)
	if err != nil {
		fmt.Printf("%v:ERRR:DB+: LoadDBs:gDB:%v\n", PRG_TAG, err)
		return err
	}
	return nil
}
