package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

func saveCacheToFile(m map[string]int64) {
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

func readCacheFromFile() map[string]int64 {
	cache := make(map[string]int64)

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
