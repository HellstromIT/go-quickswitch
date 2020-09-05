package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/viper"
)

func getDirectories(dir string) []string {
	file, err := os.Open(dir)
	var full_path_dirs []string
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	for _, v := range names {
		full_path_dirs = append(full_path_dirs, filepath.Join(dir, v))
	}
	return full_path_dirs
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/quickswitch/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	dir := viper.GetStringSlice("Directories")

	var directories []string

	for _, v := range dir {
		for _, v := range getDirectories(v) {
			//fmt.Println(v)
			directories = append(directories, v)
		}
	}

	idx, _ := fuzzyfinder.Find(directories, func(i int) string {
		return directories[i]
	})
	fmt.Println(directories[idx])
}
