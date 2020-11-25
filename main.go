package main

import (
	"flag"
	"fmt"
)

func main() {

	configfile := getConfigFile("quickswitch/quickswitch.json")

	files := readConfigFromFile(configfile)

	addPtr := flag.String("add", "", "add path to search")
	removePtr := flag.String("remove", "", "remove path from search")
	flag.Parse()

	if *addPtr != "" {
		files.addDirectory(*addPtr)
		files.saveConfigToFile(configfile)
		fmt.Printf("Directory %v added to search", *addPtr)
	} else if *removePtr != "" {
		files.removeDirectory(*removePtr)
		files.saveConfigToFile(configfile)
		fmt.Printf("Directory %v removed from search", *removePtr)
	} else {

		foundDirectories := walkDirectories(&files)

		fmt.Println(foundDirectories.getDirectory(getCwd()))
	}
}
