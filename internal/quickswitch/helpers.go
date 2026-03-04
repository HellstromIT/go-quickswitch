package quickswitch

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/HellstromIT/go-quickswitch/cmd/go-quickswitch/internal/log"
)

func findInDirectoryConf(slice []directoryConf, val string) (int, bool) {
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

func walkDir(p string, d directories, f *map[string]time.Time, depth int, maxDepth int) directories {

	d.name = p
	d.depth = depth
	var childdir []directories

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
			info, err := os.Stat(filepath.Join(p, v))
			if err != nil {
				return d
			}
			if info.IsDir() {
				childPath := filepath.Join(p, v)

				var newChild directories

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

func walkGitDir(p string, d directories, f *map[string]time.Time, depth int) directories {

	d.name = p
	d.depth = depth
	var childdir []directories

	file, err := os.Open(p)
	if err != nil {
		return d
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		return d
	}
	for _, v := range names {
		if v == ".git" {
			d.searched = true
			d.time = time.Now()
			(*f)[p] = time.Now()
			return d
		}
	}
	for _, v := range names {
		info, err := os.Stat(filepath.Join(p, v))
		if err != nil {
			return d
		}
		if !info.IsDir() {
			continue
		}
		childPath := filepath.Join(p, v)

		var newChild directories
		childdir = append(childdir, walkGitDir(childPath, newChild, f, d.depth+1))
	}

	d.child = childdir
	d.searched = true
	d.time = time.Now()
	(*f)[p] = time.Now()

	return d
}

func walk(f fileList) {
	var d directories
	d.name = "pseudo"
	var childdir []directories
	flat := make(map[string]time.Time)
	for _, dir := range f.Directories {
		log.Debug("walking directory", "path", dir.Directory, "git", dir.Git, "depth", dir.Depth)
		if dir.Git {
			var newChild directories
			childdir = append(childdir, walkGitDir(dir.Directory, newChild, &flat, 0))
		} else {
			var newChild directories
			childdir = append(childdir, walkDir(dir.Directory, newChild, &flat, 0, dir.Depth))
		}
	}

	if err := saveCacheToFile(flat); err != nil {
		log.Error("failed to save cache", "error", err)
	}
	d.child = childdir
}

func getCwd() string {
	path, err := os.Getwd()
	if err != nil {
		log.Error("failed to get current directory", "error", err)
		return "."
	}
	return path
}

// walkLive walks directories and updates both the cache map and a live list.
// The list is updated with mutex protection for hot reload support.
func walkLive(f fileList, list *[]string, mu *sync.RWMutex, seen map[string]bool) {
	flat := make(map[string]time.Time)

	for _, dir := range f.Directories {
		log.Debug("walking directory (live)", "path", dir.Directory, "git", dir.Git, "depth", dir.Depth)
		if dir.Git {
			walkGitDirLive(dir.Directory, &flat, list, mu, seen)
		} else {
			walkDirLive(dir.Directory, &flat, 0, dir.Depth, list, mu, seen)
		}
	}

	if err := saveCacheToFile(flat); err != nil {
		log.Error("failed to save cache", "error", err)
	}
}

func walkDirLive(p string, f *map[string]time.Time, depth int, maxDepth int, list *[]string, mu *sync.RWMutex, seen map[string]bool) {
	// Add this directory to the live list if not already seen
	mu.Lock()
	if !seen[p] {
		seen[p] = true
		*list = append(*list, p)
	}
	mu.Unlock()

	(*f)[p] = time.Now()

	if depth >= maxDepth {
		return
	}

	file, err := os.Open(p)
	if err != nil {
		return
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		return
	}

	for _, v := range names {
		childPath := filepath.Join(p, v)
		info, err := os.Stat(childPath)
		if err != nil {
			continue
		}
		if info.IsDir() {
			walkDirLive(childPath, f, depth+1, maxDepth, list, mu, seen)
		}
	}
}

func walkGitDirLive(p string, f *map[string]time.Time, list *[]string, mu *sync.RWMutex, seen map[string]bool) {
	file, err := os.Open(p)
	if err != nil {
		return
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		return
	}

	// Check if this is a git repo
	for _, v := range names {
		if v == ".git" {
			// Found a git repo - add it and stop recursing
			mu.Lock()
			if !seen[p] {
				seen[p] = true
				*list = append(*list, p)
			}
			mu.Unlock()
			(*f)[p] = time.Now()
			return
		}
	}

	// Not a git repo, recurse into subdirectories
	for _, v := range names {
		childPath := filepath.Join(p, v)
		info, err := os.Stat(childPath)
		if err != nil {
			continue
		}
		if info.IsDir() {
			walkGitDirLive(childPath, f, list, mu, seen)
		}
	}
}
