package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type context struct {
	configFile string
	files      FileList
}

type addCmdSub struct {
	Depth int  `short:"d" help:"How deep the crawler should traverse directory (Only relevant if git is unset). Default: 0"`
	Git   bool `short:"g" optional help:"If set crawler will look for all subdirs with a .git folder. Default: false "`
}

type addCmd struct {
	Paths     string `required arg name:"path" help:"Full Path to add." type:"path"`
	addCmdSub `cmd`
}

type rmCmd struct {
	Paths string `arg name: "path" help:"Full path to remove." type:"path"`
}

type runCmd struct {
}

var cli struct {
	Add    addCmd `cmd help:"Add Paths to configuration file."`
	Remove rmCmd  `cmd help:"Remove Paths from configuration file."`
	Run    runCmd `cmd help:"Fuzzy search directories" default:"1"`
}

func (a *addCmd) Run(ctx *context) error {
	ctx.files.addDirectory(a.Paths, a.Git, a.Depth)
	ctx.files.saveConfigToFile(ctx.configFile)
	walk(ctx.files)
	fmt.Printf("Directory %v added to search", a.Paths)
	return nil
}

func (r *rmCmd) Run(ctx *context) error {
	ctx.files.removeDirectory(r.Paths)
	ctx.files.saveConfigToFile(ctx.configFile)
	walk(ctx.files)
	fmt.Printf("Directory %v removed from search", r.Paths)
	return nil
}

func (r *runCmd) Run(ctx *context) error {
	cache := readCacheFromFile()
	go walk(ctx.files)
	fmt.Println(getDirectory(cache, getCwd()))
	return nil
}

func Cli() {
	configfile := getConfigFile("quickswitch/quickswitch.json")

	files := readConfigFromFile(configfile)

	ctx := kong.Parse(&cli)

	err := ctx.Run(&context{configFile: configfile, files: files})
	ctx.FatalIfErrorf(err)
}
