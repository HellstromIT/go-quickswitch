package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileList Holds the directories to search
type FileList struct {
	Directories []DirectoryConf
}

// DirectoryConf Holds configuration for each directory
type DirectoryConf struct {
	Directory string
	Git       bool
	Depth     int
}

type FoundDirectories struct {
	directories []string
}

type Directories struct {
	name     string
	depth    int
	searched bool
	time     time.Time
	child    []Directories
}

func (f *FileList) addDirectory(d string, git bool, depth int) {
	newDirectory := DirectoryConf{
		Directory: d,
		Git:       git,
		Depth:     depth,
	}
	_, found := findInDirectoryConf(f.Directories, d)
	if !found {
		f.Directories = append(f.Directories, newDirectory)
	}
}

func (f *FileList) removeDirectory(directory string) {
	i, found := findInDirectoryConf(f.Directories, directory)
	if !found {
		fmt.Println("Directory not found in config. Make sure you're using the exact path")
		os.Exit(1)
	}
	f.Directories = append(f.Directories[:i], f.Directories[i+1:]...)
}

func (f *FileList) createBaseConfig(filename string) {

	(*f).addDirectory(getCwd(), false, 0)

	errMkdir := os.MkdirAll(filepath.Dir(filename), 0755)
	if errMkdir != nil {
		fmt.Println("Error:", errMkdir)
		os.Exit(1)
	}

	f.saveConfigToFile(filename)
}
