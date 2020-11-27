package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

type FoundDirectories struct {
	directories []string
}

type Directories struct {
	name     string
	searched bool
	time     time.Time
	child    []Directories
}

func (f *FileList) addDirectory(d string, git bool) {
	newDirectory := DirectoryConf{
		Directory: d,
		Git:       git,
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

	(*f).addDirectory(getCwd(), false)

	errMkdir := os.MkdirAll(filepath.Dir(filename), 0755)
	if errMkdir != nil {
		fmt.Println("Error:", errMkdir)
		os.Exit(1)
	}

	f.saveConfigToFile(filename)
}

func (f *FileList) saveConfigToFile(filename string) error {
	bs, err := json.MarshalIndent(*f, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return ioutil.WriteFile(filename, bs, 0644)
}

func readConfigFromFile(filename string) FileList {
	var filelist FileList

	if _, err := os.Stat(filename); err == nil {
		bs, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonErr := json.Unmarshal(bs, &filelist)
		if jsonErr != nil {
			fmt.Println("Error:", jsonErr)
			os.Exit(1)
		}
		return filelist
	} else if os.IsNotExist(err) {
		filelist.createBaseConfig(filename)
		fmt.Printf("Creating configuration at:\n   %v\n", filename)
		fmt.Println("Configuration created. Re-run command to search")
		os.Exit(0)
	} else {
		fmt.Println("Error: Most likely .config/quickswitch is a file not a dir")
	}

	return filelist
}
