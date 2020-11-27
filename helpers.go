package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var wg sync.WaitGroup

func getConfigFile(f string) string {
	home, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error:", err)
	}
	return filepath.Join(home, f)
}

func findInDirectoryConf(slice []DirectoryConf, val string) (int, bool) {
	for i, item := range slice {
		if item.Directory == val {
			return i, true
		}
	}
	return -1, false
}

func findInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func walkDir(p string, d Directories, f *map[string]int64) Directories {

	d.name = p
	var childdir []Directories

	file, err := os.Open(p)
	if err != nil {
		return d
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		return d
	}
	for _, v := range names {
		info, err := os.Stat(p + "/" + v)
		if err != nil {
			return d
		}
		if v == ".git" {
			d.searched = true
			d.time = time.Now().Unix()
			(*f)[p] = time.Now().Unix()
			return d
		} else if !info.IsDir() {
			d.searched = true
			d.time = time.Now().Unix()
			(*f)[p] = time.Now().Unix()
			return d
		}
	}
	for _, v := range names {
		childPath := p + "/" + v

		var newChild Directories

		childdir = append(childdir, walkDir(childPath, newChild, f))

	}
	d.child = childdir
	d.searched = true
	d.time = time.Now().Unix()
	(*f)[p] = time.Now().Unix()

	return d
}

func walk(f FileList, flat map[string]int64) {
	var d Directories
	d.name = "pseudo"
	var childdir []Directories

	for _, dir := range f.Directories {
		var newChild Directories

		childdir = append(childdir, walkDir(dir.Directory, newChild, &flat))
	}

	saveCacheToFile(flat)
	d.child = childdir

	return
}

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}

	return path
}
