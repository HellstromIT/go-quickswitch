package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	configfile := getConfigFile("quickswitch/quickswitch.json")

	files := readConfigFromFile(configfile)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	//addFile := addCmd.String("path", "", "Full path to add")
	addGit := addCmd.Bool("git", false, "Should recursive search for git repo be enables")
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nSubcommands [add, remove]\n\n")
		fmt.Fprintf(os.Stderr, "List subcommand usage with:\n\n")
		fmt.Fprintf(os.Stderr, "  go-quickswitch add -h\n\n")
		fmt.Fprintf(os.Stderr, "  go-quickswitch remove -h\n\n")
	}
	flag.Parse()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add":
			addCmd.Parse(os.Args[2:])
			if len(addCmd.Args()) == 1 {
				files.addDirectory(addCmd.Args()[0], *addGit)
				files.saveConfigToFile(configfile)
				fmt.Printf("Directory %v added to search", addCmd.Args()[0])
			}
		case "remove":
			removeCmd.Parse(os.Args[2:])
			if len(removeCmd.Args()) == 1 {
				if removeCmd.Args()[0] != "" {
					files.removeDirectory(removeCmd.Args()[0])
					files.saveConfigToFile(configfile)
					fmt.Printf("Directory %v removed from search", removeCmd.Args()[0])
				}
			}
			fmt.Println("Too many arguments")
		default:
		}
	}
	foundDirectories := walkDirectories(&files)

	fmt.Println(foundDirectories.getDirectory(getCwd()))
}
