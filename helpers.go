package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func walkDir(dir string, d *FileList) {
	defer wg.Done()

	walk := func(path string, f os.FileInfo, err error) error {
		suffix := ".git"
		if f.IsDir() && strings.HasSuffix(path, suffix) {
			(*d).Lock()
			(*d).addDirectory(strings.TrimSuffix(path, suffix), false)
			(*d).Unlock()
		} else if f.IsDir() && path != dir {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error:", err)
			}
			names, err := file.Readdirnames(0)
			file.Close()
			if err != nil {
				fmt.Println("Error:", err)
			}
			for _, v := range names {
				wg.Add(1)
				go walkDir(filepath.Join(path, v), d)
			}
			return filepath.SkipDir
		}
		return nil
	}

	filepath.Walk(dir, walk)
}

func walk(f FileList, d *FileList) FileList {
	for _, dir := range f.Directories {
		wg.Add(1)
		walkDir(dir.Directory, d)
	}
	wg.Wait()
	return (*d)
}

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}

	return path
}
