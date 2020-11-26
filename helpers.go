package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func isGitDirectory(d string) bool {
	//fmt.Println(d)
	info, _ := os.Stat(d)
	if info.IsDir() {
		file, err := os.Open(d)
		if err != nil {
			fmt.Println("Error:", err)
		}

		names, err := file.Readdirnames(0)
		if err != nil {
			fmt.Println("Error:", err)
		}

		_, found := findInSlice(names, ".git")
		if !found {
			return false
		}
		return true
	}

	return false
}

func walkDir(dir string, d *FileList) {
	defer wg.Done()

	walk := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && isGitDirectory(path) {
			(*d).Lock()
			(*d).addDirectory(path, false)
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

func walkDirectories2(f FileList, s int, e int) FileList {
	var foundDir FileList

	for s < e {
		s++
		for _, dir := range f.Directories {
			info, _ := os.Stat(dir.Directory)
			if info.IsDir() && isGitDirectory(dir.Directory) {
				foundDir.addDirectory(dir.Directory, false)
			} else if info.IsDir() {
				file, err := os.Open(dir.Directory)
				if err != nil {
					fmt.Println("Error:", err)
				}
				names, err := file.Readdirnames(0)
				if err != nil {
					fmt.Println("Error:", err)
				}
				for _, v := range names {
					var recurseDir FileList
					recurseDir.addDirectory(filepath.Join(dir.Directory, v), false)
					found := walkDirectories2(recurseDir, 0, 1)
					for _, dir := range found.Directories {
						foundDir.addDirectory(dir.Directory, false)
					}
				}
			}
		}
	}
	return foundDir
}

func walkDirectories(f *FileList) FileList {

	var foundDir FileList

	for _, dir := range f.Directories {
		foundDir.addDirectory(dir.Directory, false)
		if isGitDirectory(dir.Directory) {
			continue
		}
		file, err := os.Open(dir.Directory)
		if err != nil {
			fmt.Println("Error:", err)
		}

		names, err := file.Readdirnames(0)
		if err != nil {
			fmt.Println("Error:", err)
		}

		for _, v := range names {
			var mergeDir FileList
			mergeDir.addDirectory(filepath.Join(dir.Directory, v), false)
			info, _ := os.Stat(filepath.Join(dir.Directory, v))
			if info.IsDir() && isGitDirectory(filepath.Join(dir.Directory, v)) {
				foundDir.addDirectory(filepath.Join(dir.Directory, v), false)
			} else if info.IsDir() {
				subDir := walkDirectories(&mergeDir)
				for _, d := range subDir.Directories {
					foundDir.addDirectory(d.Directory, false)
				}
			}
		}
	}
	return foundDir
}

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}

	return path
}
