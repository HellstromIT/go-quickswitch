package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	configfile := getConfigFile(".config/quickswitch.json")

	files := readConfigFromFile(configfile)

	addPtr := flag.String("add", "", "add path to search")
	flag.Parse()
	if *addPtr != "" {
		files.addDirectory(*addPtr)
		files.saveConfigToFile(configfile)
		os.Exit(0)
	}

	walkDirectories(&files)

	cwd := getCwd()

	directory := files.getDirectory(cwd)

	fmt.Println(directory)
}
