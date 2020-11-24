package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// FileList Holds the configure paths to search
type FileList struct {
	Directories []string
}

func (f *FileList) addDirectory(directory string) {
	f.Directories = append(f.Directories, directory)
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
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var filelist FileList
	jsonErr := json.Unmarshal(bs, &filelist)
	if err != nil {
		fmt.Println("Error:", jsonErr)
		os.Exit(1)
	}

	return filelist
}
