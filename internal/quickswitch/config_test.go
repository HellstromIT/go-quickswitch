package quickswitch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testPathProjects = "/home/user/projects"
	testPathWork     = "/home/user/work"
)

func TestAddDirectory(t *testing.T) {
	tests := []struct {
		name         string
		initial      []directoryConf
		addPath      string
		addGit       bool
		addDepth     int
		wantLen      int
		wantContains string
	}{
		{
			name:         "add to empty list",
			initial:      []directoryConf{},
			addPath:      testPathProjects,
			addGit:       true,
			addDepth:     0,
			wantLen:      1,
			wantContains: testPathProjects,
		},
		{
			name: "add to existing list",
			initial: []directoryConf{
				{Directory: testPathWork, Git: false, Depth: 2},
			},
			addPath:      testPathProjects,
			addGit:       true,
			addDepth:     0,
			wantLen:      2,
			wantContains: testPathProjects,
		},
		{
			name: "add duplicate is ignored",
			initial: []directoryConf{
				{Directory: testPathProjects, Git: true, Depth: 0},
			},
			addPath:      testPathProjects,
			addGit:       false,
			addDepth:     5,
			wantLen:      1,
			wantContains: testPathProjects,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileList{Directories: tt.initial}
			f.addDirectory(tt.addPath, tt.addGit, tt.addDepth)

			if len(f.Directories) != tt.wantLen {
				t.Errorf("addDirectory() resulted in %d directories, want %d", len(f.Directories), tt.wantLen)
			}

			_, found := findInDirectoryConf(f.Directories, tt.wantContains)
			if !found {
				t.Errorf("addDirectory() did not contain expected path: %s", tt.wantContains)
			}
		})
	}
}

func TestAddDirectoryPreservesFlags(t *testing.T) {
	f := &fileList{}
	f.addDirectory("/path/one", true, 0)
	f.addDirectory("/path/two", false, 5)

	if len(f.Directories) != 2 {
		t.Fatalf("expected 2 directories, got %d", len(f.Directories))
	}

	// Check first directory
	if f.Directories[0].Git != true {
		t.Error("first directory should have Git=true")
	}
	if f.Directories[0].Depth != 0 {
		t.Error("first directory should have Depth=0")
	}

	// Check second directory
	if f.Directories[1].Git != false {
		t.Error("second directory should have Git=false")
	}
	if f.Directories[1].Depth != 5 {
		t.Error("second directory should have Depth=5")
	}
}

func TestSaveAndReadConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	// Create a fileList and save it
	original := &fileList{
		Directories: []directoryConf{
			{Directory: testPathProjects, Git: true, Depth: 0},
			{Directory: testPathWork, Git: false, Depth: 3},
		},
	}

	err := original.saveConfigToFile(configFile)
	if err != nil {
		t.Fatalf("saveConfigToFile() error = %v", err)
	}

	// Verify file exists and is valid JSON
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var loaded fileList
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("config file is not valid JSON: %v", err)
	}

	// Verify contents
	if len(loaded.Directories) != 2 {
		t.Errorf("loaded config has %d directories, want 2", len(loaded.Directories))
	}

	if loaded.Directories[0].Directory != testPathProjects {
		t.Errorf("first directory = %s, want %s", loaded.Directories[0].Directory, testPathProjects)
	}

	if loaded.Directories[0].Git != true {
		t.Error("first directory Git flag should be true")
	}

	if loaded.Directories[1].Depth != 3 {
		t.Errorf("second directory Depth = %d, want 3", loaded.Directories[1].Depth)
	}
}

func TestRemoveDirectory(t *testing.T) {
	f := &fileList{
		Directories: []directoryConf{
			{Directory: testPathProjects, Git: true, Depth: 0},
			{Directory: testPathWork, Git: false, Depth: 2},
		},
	}

	// Remove existing directory
	err := f.removeDirectory(testPathWork)
	if err != nil {
		t.Errorf("removeDirectory() unexpected error: %v", err)
	}
	if len(f.Directories) != 1 {
		t.Errorf("removeDirectory() resulted in %d directories, want 1", len(f.Directories))
	}

	// Try to remove non-existent directory
	err = f.removeDirectory("/nonexistent/path")
	if err == nil {
		t.Error("removeDirectory() expected error for non-existent path")
	}
}

func TestConfigFileFormat(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	f := &fileList{
		Directories: []directoryConf{
			{Directory: "/test/path", Git: true, Depth: 2},
		},
	}

	err := f.saveConfigToFile(configFile)
	if err != nil {
		t.Fatalf("saveConfigToFile() error = %v", err)
	}

	// Read raw content and verify it's properly indented
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	content := string(data)

	// Check that it contains expected JSON structure
	if !contains(content, `"Directories"`) {
		t.Error("config should contain Directories key")
	}
	if !contains(content, `"Directory"`) {
		t.Error("config should contain Directory key")
	}
	if !contains(content, `"Git"`) {
		t.Error("config should contain Git key")
	}
	if !contains(content, `"Depth"`) {
		t.Error("config should contain Depth key")
	}
}

func TestSaveAndReadCacheFile(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.gob")

	// Create test data
	original := map[string]time.Time{
		"/home/user/project1": time.Now().Add(-1 * time.Hour),
		"/home/user/project2": time.Now(),
	}

	// Save cache
	err := saveCacheToFile(cacheFile, original)
	if err != nil {
		t.Fatalf("saveCacheToFile() error = %v", err)
	}

	// Read cache back
	loaded := readCacheFromFile(cacheFile)

	if len(loaded) != len(original) {
		t.Errorf("readCacheFromFile() returned %d entries, want %d", len(loaded), len(original))
	}

	for path := range original {
		if _, ok := loaded[path]; !ok {
			t.Errorf("readCacheFromFile() missing path: %s", path)
		}
	}
}

func TestReadCacheFromNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "nonexistent.gob")

	// Should return empty map, not error
	cache := readCacheFromFile(cacheFile)

	if cache == nil {
		t.Error("readCacheFromFile() returned nil, want empty map")
	}
	if len(cache) != 0 {
		t.Errorf("readCacheFromFile() returned %d entries, want 0", len(cache))
	}
}

func TestReadConfigFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	// Create a config file
	original := &fileList{
		Directories: []directoryConf{
			{Directory: testPathProjects, Git: true, Depth: 0},
		},
	}
	err := original.saveConfigToFile(configFile)
	if err != nil {
		t.Fatalf("saveConfigToFile() error = %v", err)
	}

	// Read it back
	result, err := readConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("readConfigFromFile() error = %v", err)
	}

	if result.Created {
		t.Error("readConfigFromFile() Created = true, want false")
	}

	if len(result.FileList.Directories) != 1 {
		t.Errorf("readConfigFromFile() returned %d directories, want 1", len(result.FileList.Directories))
	}
}

func TestReadConfigFromFileNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "subdir", "config.json")

	// Reading non-existent file should create it
	result, err := readConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("readConfigFromFile() error = %v", err)
	}

	if !result.Created {
		t.Error("readConfigFromFile() Created = false, want true")
	}

	// Verify file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("readConfigFromFile() did not create config file")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
