package quickswitch

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func (f *FileList) saveConfigToFile(filename string) error {
	bs, err := json.MarshalIndent(*f, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return ioutil.WriteFile(filename, bs, 0644)
}

func readConfigFromFile(filename string) FileList {
	var filelist FileList

	if _, err := os.Stat(filename); err == nil {
		bs, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonErr := json.Unmarshal(bs, &filelist)
		if jsonErr != nil {
			fmt.Println("Error:", jsonErr)
			os.Exit(1)
		}
		return filelist
	} else if os.IsNotExist(err) {
		filelist.createBaseConfig(filename)
		fmt.Printf("Creating configuration at:\n   %v\n", filename)
		fmt.Println("Configuration created. Re-run command to search")
		os.Exit(0)
	} else {
		fmt.Println("Error: Most likely .config/quickswitch is a file not a dir")
	}

	return filelist
}

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
