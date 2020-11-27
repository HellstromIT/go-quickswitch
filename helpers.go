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

func walkDir(p string, d Directories) Directories {

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
			return d
		} else if !info.IsDir() {
			return d
		}
	}
	for _, v := range names {
		childPath := p + "/" + v

		var newChild Directories

		childdir = append(childdir, walkDir(childPath, newChild))

	}
	d.child = childdir
	d.searched = true
	d.time = time.Now()

	return d
}

func walk(f FileList) Directories {
	var d Directories
	d.name = "pseudo"
	var childdir []Directories

	for _, dir := range f.Directories {
		var newChild Directories

		childdir = append(childdir, walkDir(dir.Directory, newChild))
	}
	d.child = childdir

	return d
}

func (s *FoundDirectories) flattenDirectories(d Directories) FoundDirectories {
	if d.name != "pseudo" {
		s.directories = append(s.directories, d.name)
	}
	for _, c := range d.child {
		if c.child != nil {
			s.flattenDirectories(c)
		} else {
			s.directories = append(s.directories, c.name)
		}
	}
	return *s
}

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}

	return path
}
