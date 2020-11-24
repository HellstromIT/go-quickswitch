package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	configfile := getConfigFile("quickswitch/quickswitch.json")

	files := readConfigFromFile(configfile)

	addPtr := flag.String("add", "", "add path to search")
	flag.Parse()
	if *addPtr != "" {
		files.addDirectory(*addPtr)
		files.saveConfigToFile(configfile)
		os.Exit(0)
	}

	walkDirectories(&files)

	directory := files.getDirectory(getCwd())

	fmt.Println(directory)
}
