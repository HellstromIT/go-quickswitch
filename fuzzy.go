package main

import (
	"github.com/ktr0731/go-fuzzyfinder"
)

func (f FileList) getDirectory(cwd string) string {
	idx, err := fuzzyfinder.Find(f.Directories, func(i int) string {
		return f.Directories[i]
	})
	if err != nil {
		return cwd
	}
	return f.Directories[idx]
}
