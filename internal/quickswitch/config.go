package quickswitch

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/HellstromIT/go-quickswitch/cmd/go-quickswitch/internal/log"
)

// ErrDirectoryNotFound is returned when a directory is not found in the config
var ErrDirectoryNotFound = errors.New("directory not found in config")

type fileList struct {
	Directories []directoryConf `json:"Directories"`
}

type directoryConf struct {
	Directory string `json:"Directory"`
	Git       bool   `json:"Git"`
	Depth     int    `json:"Depth"`
}

type directories struct {
	name     string
	depth    int
	searched bool
	time     time.Time
	child    []directories
}

func getConfigFile(f string) (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return filepath.Join(home, f), nil
}

func (f *fileList) addDirectory(d string, git bool, depth int) {
	newDirectory := directoryConf{
		Directory: d,
		Git:       git,
		Depth:     depth,
	}
	_, found := findInDirectoryConf(f.Directories, d)
	if !found {
		f.Directories = append(f.Directories, newDirectory)
		log.Debug("added directory to config", "path", d, "git", git, "depth", depth)
	} else {
		log.Debug("directory already in config", "path", d)
	}
}

func (f *fileList) removeDirectory(directory string) error {
	i, found := findInDirectoryConf(f.Directories, directory)
	if !found {
		return fmt.Errorf("%w: %s", ErrDirectoryNotFound, directory)
	}
	f.Directories = append(f.Directories[:i], f.Directories[i+1:]...)
	log.Debug("removed directory from config", "path", directory)
	return nil
}

func (f *fileList) createBaseConfig(filename string) error {
	f.addDirectory(getCwd(), false, 0)

	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return f.saveConfigToFile(filename)
}

func (f *fileList) saveConfigToFile(filename string) error {
	bs, err := json.MarshalIndent(*f, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(filename, bs, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	log.Debug("saved config to file", "path", filename)
	return nil
}

// ConfigResult holds the result of reading a config file
type ConfigResult struct {
	FileList fileList
	Created  bool // true if a new config was created
}

func readConfigFromFile(filename string) (ConfigResult, error) {
	var result ConfigResult

	stat, err := os.Stat(filename)
	if err == nil {
		// File exists
		if stat.IsDir() {
			return result, fmt.Errorf("config path is a directory: %s", filename)
		}

		bs, err := os.ReadFile(filename)
		if err != nil {
			return result, fmt.Errorf("failed to read config file: %w", err)
		}

		err = json.Unmarshal(bs, &result.FileList)
		if err != nil {
			return result, fmt.Errorf("failed to parse config file: %w", err)
		}
		log.Debug("loaded config from file", "path", filename, "directories", len(result.FileList.Directories))
		return result, nil
	}

	if os.IsNotExist(err) {
		// Create new config
		if err := result.FileList.createBaseConfig(filename); err != nil {
			return result, err
		}
		result.Created = true
		log.Debug("created new config file", "path", filename)
		return result, nil
	}

	return result, fmt.Errorf("failed to stat config file: %w", err)
}

func saveCacheToFile(m map[string]time.Time) error {
	cachePath, err := getConfigFile("quickswitch/cache.json")
	if err != nil {
		return err
	}

	file, err := os.Create(cachePath)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	e := gob.NewEncoder(file)
	if err := e.Encode(m); err != nil {
		return fmt.Errorf("failed to encode cache: %w", err)
	}
	log.Debug("saved cache to file", "path", cachePath, "entries", len(m))
	return nil
}

func readCacheFromFile() map[string]time.Time {
	cache := make(map[string]time.Time)

	cachePath, err := getConfigFile("quickswitch/cache.json")
	if err != nil {
		log.Error("failed to get cache path", "error", err)
		return cache
	}

	file, err := os.Open(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Error("failed to open cache file", "error", err)
		}
		return cache
	}
	defer file.Close()

	d := gob.NewDecoder(file)
	if err := d.Decode(&cache); err != nil {
		log.Error("failed to decode cache", "error", err)
		return cache
	}
	log.Debug("loaded cache from file", "path", cachePath, "entries", len(cache))
	return cache
}
