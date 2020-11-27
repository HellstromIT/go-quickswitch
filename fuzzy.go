package main

import (
	"github.com/ktr0731/go-fuzzyfinder"
)

func (f FoundDirectories) getDirectory(cwd string) string {
	idx, err := fuzzyfinder.Find(f.directories, func(i int) string {
		return f.directories[i]
	})
	if err != nil {
		return cwd
	}
	return f.directories[idx]
}
