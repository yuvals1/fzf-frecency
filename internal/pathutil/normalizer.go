package pathutil

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

// PathNormalizer handles path normalization operations
type PathNormalizer struct {
    currentDir string
    homeDir    string
}

// NewPathNormalizer creates a new PathNormalizer
func NewPathNormalizer() (*PathNormalizer, error) {
    currentDir, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("getting current directory: %w", err)
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("getting home directory: %w", err)
    }

    return &PathNormalizer{
        currentDir: currentDir,
        homeDir:    homeDir,
    }, nil
}

// stripANSI removes ANSI escape sequences efficiently
func stripANSI(s string) string {
    var b strings.Builder
    b.Grow(len(s))
    inEscape := false

    for i := 0; i < len(s); i++ {
        if s[i] == '\x1b' {
            inEscape = true
            continue
        }
        if inEscape {
            if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
                inEscape = false
            }
            continue
        }
        b.WriteByte(s[i])
    }
    return b.String()
}

// ToStoragePath converts a path to absolute format for storage
func (pn *PathNormalizer) ToStoragePath(path string) (string, error) {
    // Clean and normalize the path
    path = stripANSI(strings.TrimSpace(path))
    
    // If it's already absolute, just clean it
    if filepath.IsAbs(path) {
        return filepath.Clean(path), nil
    }

    // Handle home directory notation
    if strings.HasPrefix(path, "~/") {
        return filepath.Clean(filepath.Join(pn.homeDir, path[2:])), nil
    }

    // Convert relative to absolute using current directory
    absPath := filepath.Clean(filepath.Join(pn.currentDir, path))
    return absPath, nil
}

// ToDisplayPath converts a storage path to a display format
func (pn *PathNormalizer) ToDisplayPath(storagePath string) (string, error) {
    // Clean storage path
    storagePath = filepath.Clean(storagePath)
    
    // Try to make it relative to current directory first
    relToCurrent, err := filepath.Rel(pn.currentDir, storagePath)
    if err == nil && !strings.HasPrefix(relToCurrent, "..") {
        return relToCurrent, nil
    }

    // If it's under home directory, use ~ notation
    if strings.HasPrefix(storagePath, pn.homeDir) {
        rel, err := filepath.Rel(pn.homeDir, storagePath)
        if err == nil {
            return "~/" + rel, nil
        }
    }

    // Return absolute path as fallback
    return storagePath, nil
}

// IsOutsideCurrentDir checks if a path is outside the current directory
func (pn *PathNormalizer) IsOutsideCurrentDir(path string) bool {
    if filepath.IsAbs(path) {
        return !strings.HasPrefix(filepath.Clean(path), filepath.Clean(pn.currentDir))
    }
    return strings.HasPrefix(path, "..")
}

// GetCurrentDir returns the current working directory
func (pn *PathNormalizer) GetCurrentDir() string {
    return pn.currentDir
}

// GetHomeDir returns the user's home directory
func (pn *PathNormalizer) GetHomeDir() string {
    return pn.homeDir
}
