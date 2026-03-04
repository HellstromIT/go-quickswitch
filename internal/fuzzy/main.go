package fuzzy

import (
	"sync"

	"github.com/ktr0731/go-fuzzyfinder"
)

// GetDirectoryLive spawns a fuzzy finder with hot reload support.
// The list can be updated externally while the finder is running.
// Callers must hold mu.Lock() when modifying the list.
func GetDirectoryLive(list *[]string, mu *sync.RWMutex, cwd string) string {
	idx, err := fuzzyfinder.Find(
		list,
		func(i int) string {
			return (*list)[i]
		},
		fuzzyfinder.WithHotReloadLock(mu.RLocker()),
	)
	if err != nil {
		return cwd
	}
	return (*list)[idx]
}
