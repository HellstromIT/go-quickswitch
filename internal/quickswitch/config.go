package quickswitch

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type fileList struct {
	directories []directoryConf
}

type directoryConf struct {
	Directory string
	Git       bool
	Depth     int
}

type directories struct {
	name     string
	depth    int
	searched bool
	time     time.Time
	child    []directories
}

func getConfigFile(f string) string {
	home, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error:", err)
	}
	return filepath.Join(home, f)
}

func (f *fileList) addDirectory(d string, git bool, depth int) {
	newDirectory := directoryConf{
		Directory: d,
		Git:       git,
		Depth:     depth,
	}
	_, found := findInDirectoryConf(f.directories, d)
	if !found {
		f.directories = append(f.directories, newDirectory)
	}
}

func (f *fileList) removeDirectory(directory string) {
	i, found := findInDirectoryConf(f.directories, directory)
	if !found {
		fmt.Println("Directory not found in config. Make sure you're using the exact path")
		os.Exit(1)
	}
	f.directories = append(f.directories[:i], f.directories[i+1:]...)
}

func (f *fileList) createBaseConfig(filename string) {

	(*f).addDirectory(getCwd(), false, 0)

	errMkdir := os.MkdirAll(filepath.Dir(filename), 0755)
	if errMkdir != nil {
		fmt.Println("Error:", errMkdir)
		os.Exit(1)
	}

	f.saveConfigToFile(filename)
}

func (f *fileList) saveConfigToFile(filename string) error {
	bs, err := json.MarshalIndent(*f, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return ioutil.WriteFile(filename, bs, 0644)
}

func readConfigFromFile(filename string) fileList {
	var fileList fileList

	if _, err := os.Stat(filename); err == nil {
		bs, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonErr := json.Unmarshal(bs, &fileList)
		if jsonErr != nil {
			fmt.Println("Error:", jsonErr)
			os.Exit(1)
		}
		return fileList
	} else if os.IsNotExist(err) {
		fileList.createBaseConfig(filename)
		fmt.Printf("Creating configuration at:\n   %v\n", filename)
		fmt.Println("Configuration created. Re-run command to search")
		os.Exit(0)
	} else {
		fmt.Println("Error: Most likely .config/quickswitch is a file not a dir")
	}

	return fileList
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
