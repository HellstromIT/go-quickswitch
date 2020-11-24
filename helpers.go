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

func walkDirectories(f *FileList) {
	for _, dir := range f.Directories {
		file, err := os.Open(dir)
		if err != nil {
			fmt.Println("Error:", err)
		}

		names, err := file.Readdirnames(0)
		if err != nil {
			fmt.Println("Error:", err)
		}

		for _, v := range names {
			info, _ := os.Stat(filepath.Join(dir, v))
			if info.IsDir() {
				f.addDirectory(filepath.Join(dir, v))
			}
		}
	}
}

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}

	return path
}
