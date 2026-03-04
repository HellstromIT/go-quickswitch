package quickswitch

import (
	"fmt"
	"sync"

	"github.com/HellstromIT/go-quickswitch/cmd/go-quickswitch/internal/fuzzy"
	"github.com/alecthomas/kong"
)

type context struct {
	version    string
	configFile string
	files      fileList
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

type versionCmd struct {
}

type runCmd struct {
}

var cli struct {
	Add     addCmd     `cmd help:"Add Paths to configuration file."`
	Remove  rmCmd      `cmd help:"Remove Paths from configuration file."`
	Run     runCmd     `cmd help:"Fuzzy search directories" default:"1"`
	Version versionCmd `cmd help:"Print version."`
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

func (v *versionCmd) Run(ctx *context) error {
	fmt.Println(ctx.version)
	return nil
}

func (r *runCmd) Run(ctx *context) error {
	cache := readCacheFromFile()

	// Initialize list from cache
	var list []string
	seen := make(map[string]bool)
	for path := range cache {
		list = append(list, path)
		seen[path] = true
	}

	var mu sync.RWMutex
	var wg sync.WaitGroup

	// Start walking directories in background with hot reload
	wg.Add(1)
	go func() {
		defer wg.Done()
		walkLive(ctx.files, &list, &mu, seen)
	}()

	// Show fuzzy finder with hot reload support
	fmt.Println(fuzzy.GetDirectoryLive(&list, &mu, getCwd()))
	wg.Wait()
	return nil
}

// Cli func
func Cli(v string) {
	configfile := getConfigFile("quickswitch/quickswitch.json")

	files := readConfigFromFile(configfile)

	ctx := kong.Parse(&cli)

	err := ctx.Run(&context{version: v, configFile: configfile, files: files})
	ctx.FatalIfErrorf(err)
}
