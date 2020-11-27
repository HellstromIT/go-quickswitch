package main

import (
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
)

func getDirectory(m map[string]time.Time, cwd string) string {
	var list []string

	for i := range m {
		list = append(list, i)
	}
	idx, err := fuzzyfinder.Find(list, func(i int) string {
		return list[i]
	})
	if err != nil {
		return cwd
	}
	return list[idx]
}
