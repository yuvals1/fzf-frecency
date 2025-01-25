package finder

import (
    "io/fs"
    "os"
    "path/filepath"
    "strings"
)

// FileFinder handles file discovery and filtering
type FileFinder struct {
    excludePatterns []string
    includeHidden   bool
}

// NewFileFinder creates a new FileFinder with default settings
func NewFileFinder() *FileFinder {
    return &FileFinder{
        excludePatterns: []string{".git", ".mypy_cache"},
        includeHidden:   true,
    }
}

// WithExcludePatterns adds patterns to exclude from search
func (f *FileFinder) WithExcludePatterns(patterns []string) *FileFinder {
    f.excludePatterns = append(f.excludePatterns, patterns...)
    return f
}

// WithHidden sets whether to include hidden files
func (f *FileFinder) WithHidden(include bool) *FileFinder {
    f.includeHidden = include
    return f
}

// shouldExclude checks if a path should be excluded based on patterns
func (f *FileFinder) shouldExclude(path string) bool {
    // Always exclude "." and ".."
    if path == "." || path == ".." {
        return true
    }

    // Check against exclude patterns
    for _, pattern := range f.excludePatterns {
        if strings.Contains(path, pattern) {
            return true
        }
    }

    // Handle hidden files
    if !f.includeHidden && strings.HasPrefix(filepath.Base(path), ".") {
        return true
    }

    return false
}

// FindFiles returns a channel of found files
func (f *FileFinder) FindFiles(root string) (<-chan string, error) {
    filesChan := make(chan string)

    // Verify root exists
    _, err := os.Stat(root)
    if err != nil {
        return nil, err
    }

    go func() {
        defer close(filesChan)

        filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
            if err != nil {
                return err
            }

            // Skip directories and excluded paths
            if d.IsDir() || f.shouldExclude(path) {
                return nil
            }

            // Make path relative to root
            relPath, err := filepath.Rel(root, path)
            if err != nil {
                return err
            }

            filesChan <- relPath
            return nil
        })
    }()

    return filesChan, nil
}
