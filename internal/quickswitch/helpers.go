package quickswitch

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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

func walkDir(p string, d Directories, f *map[string]time.Time, depth int, maxDepth int) Directories {

	d.name = p
	d.depth = depth
	var childdir []Directories

	if d.depth < maxDepth {
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
			if info.IsDir() {
				childPath := p + "/" + v

				var newChild Directories

				childdir = append(childdir, walkDir(childPath, newChild, f, d.depth+1, maxDepth))
			}
		}
	}
	d.child = childdir
	d.searched = true
	d.time = time.Now()
	(*f)[p] = time.Now()

	return d
}

func walkGitDir(p string, d Directories, f *map[string]time.Time, depth int) Directories {

	d.name = p
	d.depth = depth
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
			d.searched = true
			d.time = time.Now()
			(*f)[p] = time.Now()
			return d
		} else if !info.IsDir() {
			d.searched = true
			d.time = time.Now()
			(*f)[p] = time.Now()
			return d
		}
	}
	for _, v := range names {
		childPath := p + "/" + v

		var newChild Directories

		childdir = append(childdir, walkGitDir(childPath, newChild, f, d.depth+1))

	}
	d.child = childdir
	d.searched = true
	d.time = time.Now()
	(*f)[p] = time.Now()

	return d
}

func walk(f FileList) {
	var d Directories
	d.name = "pseudo"
	var childdir []Directories
	flat := make(map[string]time.Time)
	for _, dir := range f.Directories {
		if dir.Git {
			var newChild Directories

			childdir = append(childdir, walkGitDir(dir.Directory, newChild, &flat, 0))
		} else {
			var newChild Directories

			childdir = append(childdir, walkDir(dir.Directory, newChild, &flat, 0, dir.Depth))
		}
	}

	saveCacheToFile(flat)
	d.child = childdir

	return
}

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}

	return path
}