package patch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Patch represents a single shell patch file to be applied.
type Patch struct {
	Name    string
	Path    string
	Content string
}

// LoadPatches reads all .sh files from the given directory,
// returning them sorted lexicographically (e.g. 001_init.sh, 002_setup.sh).
func LoadPatches(dir string) ([]Patch, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading patch directory %q: %w", dir, err)
	}

	var patches []Patch
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sh") {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("reading patch file %q: %w", fullPath, err)
		}

		patches = append(patches, Patch{
			Name:    entry.Name(),
			Path:    fullPath,
			Content: string(data),
		})
	}

	sort.Slice(patches, func(i, j int) bool {
		return patches[i].Name < patches[j].Name
	})

	return patches, nil
}
