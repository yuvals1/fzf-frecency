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
    useColor        bool
}

// NewFileFinder creates a new FileFinder with default settings
func NewFileFinder() *FileFinder {
    return &FileFinder{
        excludePatterns: []string{"*.mypy*", "*.git*"},
        includeHidden:   true,
        useColor:        true,
    }
}

// shouldExclude checks if a path should be excluded based on patterns
func (f *FileFinder) shouldExclude(path string) bool {
    // Always exclude "." and ".."
    if path == "." || path == ".." {
        return true
    }

    // Check against exclude patterns
    for _, pattern := range f.excludePatterns {
        if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
            return true
        }
        if strings.Contains(path, strings.TrimSuffix(pattern, "*")) {
            return true
        }
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

            // Skip excluded paths
            if f.shouldExclude(path) {
                if d.IsDir() {
                    return filepath.SkipDir
                }
                return nil
            }

            // Skip directories
            if d.IsDir() {
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
