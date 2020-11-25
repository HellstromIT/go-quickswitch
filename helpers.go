package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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
