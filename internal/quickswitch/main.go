package quickswitch

import (
	"fmt"
	"os"
	"sync"

	"github.com/HellstromIT/go-quickswitch/cmd/go-quickswitch/internal/fuzzy"
	"github.com/HellstromIT/go-quickswitch/cmd/go-quickswitch/internal/log"
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
	Debug   bool       `short:"D" help:"Enable debug logging"`
	Add     addCmd     `cmd help:"Add Paths to configuration file."`
	Remove  rmCmd      `cmd help:"Remove Paths from configuration file."`
	Run     runCmd     `cmd help:"Fuzzy search directories" default:"1"`
	Version versionCmd `cmd help:"Print version."`
}

func (a *addCmd) Run(ctx *context) error {
	ctx.files.addDirectory(a.Paths, a.Git, a.Depth)
	if err := ctx.files.saveConfigToFile(ctx.configFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	walk(ctx.files)
	fmt.Printf("Directory %v added to search\n", a.Paths)
	return nil
}

func (r *rmCmd) Run(ctx *context) error {
	if err := ctx.files.removeDirectory(r.Paths); err != nil {
		return err
	}
	if err := ctx.files.saveConfigToFile(ctx.configFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	walk(ctx.files)
	fmt.Printf("Directory %v removed from search\n", r.Paths)
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
	ctx := kong.Parse(&cli)

	// Enable debug logging if requested
	if cli.Debug {
		log.EnableDebug()
		log.Debug("debug logging enabled")
	}

	configfile, err := getConfigFile("quickswitch/quickswitch.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	result, err := readConfigFromFile(configfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if result.Created {
		fmt.Printf("Created configuration at:\n   %v\n", configfile)
		fmt.Println("Configuration created. Re-run command to search")
		os.Exit(0)
	}

	err = ctx.Run(&context{version: v, configFile: configfile, files: result.FileList})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
