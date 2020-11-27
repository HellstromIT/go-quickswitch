package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	configfile := getConfigFile("quickswitch/quickswitch2.json")

	files := readConfigFromFile(configfile)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDepth := addCmd.Int("depth", 0, "Depth to search for directories")
	addGit := addCmd.Bool("git", false, "Should recursive search for git repo be enabled")
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
				files.addDirectory(addCmd.Args()[0], *addGit, *addDepth)
				files.saveConfigToFile(configfile)
				walk(files)
				fmt.Printf("Directory %v added to search", addCmd.Args()[0])
				os.Exit(0)
			}
		case "remove":
			removeCmd.Parse(os.Args[2:])
			if len(removeCmd.Args()) == 1 {
				if removeCmd.Args()[0] != "" {
					files.removeDirectory(removeCmd.Args()[0])
					files.saveConfigToFile(configfile)
					walk(files)
					fmt.Printf("Directory %v removed from search", removeCmd.Args()[0])
					os.Exit(0)
				}
			}
			fmt.Println("Too many arguments")
		default:
		}
	}

	cache := readCacheFromFile()
	go walk(files)
	fmt.Println(getDirectory(cache, getCwd()))
}
