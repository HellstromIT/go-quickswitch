package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

func saveCacheToFile(m map[string]time.Time) {
	file, err := os.Create(getConfigFile("quickswitch/cache.json"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer file.Close()

	e := gob.NewEncoder(file)

	err = e.Encode(m)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return
}

func readCacheFromFile() map[string]time.Time {
	cache := make(map[string]time.Time)

	file, err := os.Open(getConfigFile("quickswitch/cache.json"))
	if err != nil {
		fmt.Println("Error:", err)
		return cache
	}
	defer file.Close()

	d := gob.NewDecoder(file)

	err = d.Decode(&cache)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return cache
}
