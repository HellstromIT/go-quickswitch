package quickswitch

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// mustMkdirAll creates a directory and fails the test if it errors
func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create test directory %s: %v", path, err)
	}
}

func TestFindInSlice(t *testing.T) {
	tests := []struct {
		name      string
		slice     []string
		val       string
		wantIndex int
		wantFound bool
	}{
		{
			name:      "find existing element",
			slice:     []string{"a", "b", "c"},
			val:       "b",
			wantIndex: 1,
			wantFound: true,
		},
		{
			name:      "find first element",
			slice:     []string{"a", "b", "c"},
			val:       "a",
			wantIndex: 0,
			wantFound: true,
		},
		{
			name:      "find last element",
			slice:     []string{"a", "b", "c"},
			val:       "c",
			wantIndex: 2,
			wantFound: true,
		},
		{
			name:      "element not found",
			slice:     []string{"a", "b", "c"},
			val:       "d",
			wantIndex: -1,
			wantFound: false,
		},
		{
			name:      "empty slice",
			slice:     []string{},
			val:       "a",
			wantIndex: -1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotFound := findInSlice(tt.slice, tt.val)
			if gotIndex != tt.wantIndex {
				t.Errorf("findInSlice() index = %v, want %v", gotIndex, tt.wantIndex)
			}
			if gotFound != tt.wantFound {
				t.Errorf("findInSlice() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestFindInDirectoryConf(t *testing.T) {
	tests := []struct {
		name      string
		slice     []directoryConf
		val       string
		wantIndex int
		wantFound bool
	}{
		{
			name: "find existing directory",
			slice: []directoryConf{
				{Directory: testPathProjects, Git: true, Depth: 0},
				{Directory: testPathWork, Git: false, Depth: 2},
			},
			val:       testPathWork,
			wantIndex: 1,
			wantFound: true,
		},
		{
			name: "directory not found",
			slice: []directoryConf{
				{Directory: testPathProjects, Git: true, Depth: 0},
			},
			val:       "/home/user/other",
			wantIndex: -1,
			wantFound: false,
		},
		{
			name:      "empty slice",
			slice:     []directoryConf{},
			val:       testPathProjects,
			wantIndex: -1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotFound := findInDirectoryConf(tt.slice, tt.val)
			if gotIndex != tt.wantIndex {
				t.Errorf("findInDirectoryConf() index = %v, want %v", gotIndex, tt.wantIndex)
			}
			if gotFound != tt.wantFound {
				t.Errorf("findInDirectoryConf() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestWalkDir(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create: tmpDir/a/b/c
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	mustMkdirAll(t, nestedDir)

	// Create a file to ensure it's not included
	testFile := filepath.Join(tmpDir, "a", "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name           string
		maxDepth       int
		wantDirCount   int
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:         "depth 0 returns only root",
			maxDepth:     0,
			wantDirCount: 1,
			wantContains: []string{tmpDir},
		},
		{
			name:         "depth 1 returns root and first level",
			maxDepth:     1,
			wantDirCount: 2,
			wantContains: []string{tmpDir, filepath.Join(tmpDir, "a")},
		},
		{
			name:         "depth 3 returns all directories",
			maxDepth:     3,
			wantDirCount: 4,
			wantContains: []string{tmpDir, filepath.Join(tmpDir, "a"), filepath.Join(tmpDir, "a", "b"), filepath.Join(tmpDir, "a", "b", "c")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flat := make(map[string]time.Time)
			var d directories
			walkDir(tmpDir, d, &flat, 0, tt.maxDepth)

			if len(flat) != tt.wantDirCount {
				t.Errorf("walkDir() found %d directories, want %d", len(flat), tt.wantDirCount)
			}

			for _, dir := range tt.wantContains {
				if _, ok := flat[dir]; !ok {
					t.Errorf("walkDir() missing expected directory: %s", dir)
				}
			}
		})
	}
}

func TestWalkGitDir(t *testing.T) {
	// Create a temporary directory structure with .git folders
	tmpDir := t.TempDir()

	// Create: tmpDir/repo1/.git (git repo)
	repo1 := filepath.Join(tmpDir, "repo1")
	mustMkdirAll(t, filepath.Join(repo1, ".git"))

	// Create: tmpDir/repo1/subdir (should not be included, parent is git repo)
	mustMkdirAll(t, filepath.Join(repo1, "subdir"))

	// Create: tmpDir/notrepo/nested/repo2/.git
	repo2 := filepath.Join(tmpDir, "notrepo", "nested", "repo2")
	mustMkdirAll(t, filepath.Join(repo2, ".git"))

	flat := make(map[string]time.Time)
	var d directories
	walkGitDir(tmpDir, d, &flat, 0)

	// Should find repo1 and repo2, but not tmpDir, notrepo, nested, or subdir
	expectedRepos := []string{repo1, repo2}
	for _, repo := range expectedRepos {
		if _, ok := flat[repo]; !ok {
			t.Errorf("walkGitDir() missing expected git repo: %s", repo)
		}
	}

	// repo1/subdir should NOT be in the results (walkGitDir stops at .git)
	subdir := filepath.Join(repo1, "subdir")
	if _, ok := flat[subdir]; ok {
		t.Errorf("walkGitDir() should not include subdirs of git repos: %s", subdir)
	}
}

func TestWalkDirLive(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create: tmpDir/a/b
	nestedDir := filepath.Join(tmpDir, "a", "b")
	mustMkdirAll(t, nestedDir)

	var list []string
	var mu sync.RWMutex
	seen := make(map[string]bool)
	flat := make(map[string]time.Time)

	walkDirLive(tmpDir, &flat, 0, 2, &list, &mu, seen)

	// Should find tmpDir, tmpDir/a, tmpDir/a/b
	expectedDirs := []string{tmpDir, filepath.Join(tmpDir, "a"), filepath.Join(tmpDir, "a", "b")}

	if len(list) != len(expectedDirs) {
		t.Errorf("walkDirLive() found %d directories, want %d", len(list), len(expectedDirs))
	}

	for _, dir := range expectedDirs {
		if !seen[dir] {
			t.Errorf("walkDirLive() missing expected directory: %s", dir)
		}
	}

	// Verify cache map is also populated
	for _, dir := range expectedDirs {
		if _, ok := flat[dir]; !ok {
			t.Errorf("walkDirLive() cache missing expected directory: %s", dir)
		}
	}
}

func TestWalkGitDirLive(t *testing.T) {
	// Create a temporary directory structure with .git folders
	tmpDir := t.TempDir()

	// Create: tmpDir/repo1/.git
	repo1 := filepath.Join(tmpDir, "repo1")
	mustMkdirAll(t, filepath.Join(repo1, ".git"))

	// Create: tmpDir/notrepo/repo2/.git
	repo2 := filepath.Join(tmpDir, "notrepo", "repo2")
	mustMkdirAll(t, filepath.Join(repo2, ".git"))

	var list []string
	var mu sync.RWMutex
	seen := make(map[string]bool)
	flat := make(map[string]time.Time)

	walkGitDirLive(tmpDir, &flat, &list, &mu, seen)

	// Should find repo1 and repo2
	expectedRepos := []string{repo1, repo2}

	if len(list) != len(expectedRepos) {
		t.Errorf("walkGitDirLive() found %d repos, want %d", len(list), len(expectedRepos))
	}

	for _, repo := range expectedRepos {
		if !seen[repo] {
			t.Errorf("walkGitDirLive() missing expected repo: %s", repo)
		}
	}
}

func TestWalkLiveSkipsDuplicates(t *testing.T) {
	tmpDir := t.TempDir()

	var list []string
	var mu sync.RWMutex
	seen := make(map[string]bool)
	flat := make(map[string]time.Time)

	// Pre-populate seen with tmpDir
	seen[tmpDir] = true
	list = append(list, tmpDir)

	walkDirLive(tmpDir, &flat, 0, 0, &list, &mu, seen)

	// Should not add tmpDir again
	if len(list) != 1 {
		t.Errorf("walkDirLive() should not add duplicates, got %d items", len(list))
	}
}

func TestWalk(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.gob")

	// Create directory structure
	projectDir := filepath.Join(tmpDir, "projects")
	mustMkdirAll(t, filepath.Join(projectDir, "subdir1"))
	mustMkdirAll(t, filepath.Join(projectDir, "subdir2"))

	// Create fileList config
	files := fileList{
		Directories: []directoryConf{
			{Directory: projectDir, Git: false, Depth: 1},
		},
	}

	// Run walk
	walk(files, cacheFile)

	// Verify cache was created
	cache := readCacheFromFile(cacheFile)
	if len(cache) == 0 {
		t.Error("walk() did not create cache entries")
	}

	// Should have projectDir, subdir1, subdir2
	expectedDirs := []string{
		projectDir,
		filepath.Join(projectDir, "subdir1"),
		filepath.Join(projectDir, "subdir2"),
	}
	for _, dir := range expectedDirs {
		if _, ok := cache[dir]; !ok {
			t.Errorf("walk() cache missing expected directory: %s", dir)
		}
	}
}

func TestWalkWithGitMode(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.gob")

	// Create git repo structure
	repo1 := filepath.Join(tmpDir, "projects", "repo1")
	mustMkdirAll(t, filepath.Join(repo1, ".git"))
	// Add a subdir inside repo1 - should NOT be in cache (stops at .git)
	mustMkdirAll(t, filepath.Join(repo1, "src"))

	repo2 := filepath.Join(tmpDir, "projects", "repo2")
	mustMkdirAll(t, filepath.Join(repo2, ".git"))

	// Create fileList config with git mode
	files := fileList{
		Directories: []directoryConf{
			{Directory: filepath.Join(tmpDir, "projects"), Git: true, Depth: 0},
		},
	}

	// Run walk
	walk(files, cacheFile)

	// Verify cache
	cache := readCacheFromFile(cacheFile)

	// Should have repo1 and repo2
	if _, ok := cache[repo1]; !ok {
		t.Errorf("walk() cache missing git repo: %s", repo1)
	}
	if _, ok := cache[repo2]; !ok {
		t.Errorf("walk() cache missing git repo: %s", repo2)
	}

	// repo1/src should NOT be in cache (walkGitDir stops at .git)
	srcDir := filepath.Join(repo1, "src")
	if _, ok := cache[srcDir]; ok {
		t.Errorf("walk() should not recurse into git repo subdirs: %s", srcDir)
	}
}

func TestWalkLive(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.gob")

	// Create directory structure
	projectDir := filepath.Join(tmpDir, "projects")
	mustMkdirAll(t, filepath.Join(projectDir, "subdir1"))

	// Create fileList config
	files := fileList{
		Directories: []directoryConf{
			{Directory: projectDir, Git: false, Depth: 1},
		},
	}

	// Setup live walk state
	var list []string
	var mu sync.RWMutex
	seen := make(map[string]bool)

	// Run walkLive
	walkLive(files, &list, &mu, seen, cacheFile)

	// Verify list was populated
	if len(list) == 0 {
		t.Error("walkLive() did not populate list")
	}

	// Verify cache was created
	cache := readCacheFromFile(cacheFile)
	if len(cache) == 0 {
		t.Error("walkLive() did not create cache entries")
	}

	// List and cache should have same entries
	if len(list) != len(cache) {
		t.Errorf("walkLive() list has %d entries, cache has %d", len(list), len(cache))
	}
}

func TestWalkLiveWithGitMode(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.gob")

	// Create git repo
	repo := filepath.Join(tmpDir, "myrepo")
	mustMkdirAll(t, filepath.Join(repo, ".git"))

	// Create fileList config with git mode
	files := fileList{
		Directories: []directoryConf{
			{Directory: tmpDir, Git: true, Depth: 0},
		},
	}

	// Setup live walk state
	var list []string
	var mu sync.RWMutex
	seen := make(map[string]bool)

	// Run walkLive
	walkLive(files, &list, &mu, seen, cacheFile)

	// Should find the repo
	found := false
	for _, path := range list {
		if path == repo {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("walkLive() did not find git repo: %s", repo)
	}
}
