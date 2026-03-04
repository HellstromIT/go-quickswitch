package quickswitch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
				{Directory: "/home/user/projects", Git: true, Depth: 0},
				{Directory: "/home/user/work", Git: false, Depth: 2},
			},
			val:       "/home/user/work",
			wantIndex: 1,
			wantFound: true,
		},
		{
			name: "directory not found",
			slice: []directoryConf{
				{Directory: "/home/user/projects", Git: true, Depth: 0},
			},
			val:       "/home/user/other",
			wantIndex: -1,
			wantFound: false,
		},
		{
			name:      "empty slice",
			slice:     []directoryConf{},
			val:       "/home/user/projects",
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
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("failed to create test directories: %v", err)
	}

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
	if err := os.MkdirAll(filepath.Join(repo1, ".git"), 0755); err != nil {
		t.Fatalf("failed to create test directories: %v", err)
	}

	// Create: tmpDir/repo1/subdir (should not be included, parent is git repo)
	if err := os.MkdirAll(filepath.Join(repo1, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create test directories: %v", err)
	}

	// Create: tmpDir/notrepo/nested/repo2/.git
	repo2 := filepath.Join(tmpDir, "notrepo", "nested", "repo2")
	if err := os.MkdirAll(filepath.Join(repo2, ".git"), 0755); err != nil {
		t.Fatalf("failed to create test directories: %v", err)
	}

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
