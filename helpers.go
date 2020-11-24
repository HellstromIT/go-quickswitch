package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func getConfigFile(f string) string {
	home, err := homedir.Dir()
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
			f.addDirectory(filepath.Join(dir, v))
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
