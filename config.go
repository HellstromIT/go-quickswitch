package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileList Holds the configure paths to search
type FileList struct {
	Directories []string
}

func (f *FileList) addDirectory(directory string) {
	f.Directories = append(f.Directories, directory)
}

func (f *FileList) createBaseConfig(filename string) {

	(*f).addDirectory(filepath.Dir(filename))

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
		fmt.Println("Error", err)
		os.Exit(1)
	}
	return ioutil.WriteFile(filename, bs, 0644)
}

func readConfigFromFile(filename string) FileList {
	var filelist FileList

	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		filelist.createBaseConfig(filename)
	}

	jsonErr := json.Unmarshal(bs, &filelist)
	if jsonErr != nil {
		fmt.Println("Error:", jsonErr)
		fmt.Println("This could be because of first run. Run go-quickswitch -add=/path/to/dir to add you're first directory")
		os.Exit(1)
	}

	return filelist
}
