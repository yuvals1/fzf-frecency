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

// ToStoragePath converts a path to absolute format for storage
func (pn *PathNormalizer) ToStoragePath(path string) (string, error) {
    // Clean and trim the path
    path = strings.TrimSpace(filepath.Clean(path))
    
    // If it's already absolute, just clean it
    if filepath.IsAbs(path) {
        return filepath.Clean(path), nil
    }
    
    // Convert relative path to absolute using current directory
    return filepath.Abs(filepath.Join(pn.currentDir, path))
}

// ToDisplayPath converts a storage path to a display format
func (pn *PathNormalizer) ToDisplayPath(storagePath string) (string, error) {
    // Clean the storage path
    storagePath = filepath.Clean(storagePath)
    
    // Try to make it relative to current directory first
    if rel, err := filepath.Rel(pn.currentDir, storagePath); err == nil && !strings.HasPrefix(rel, "..") {
        return rel, nil
    }
    
    // If it's under home directory, use ~ notation
    if strings.HasPrefix(storagePath, pn.homeDir) {
        rel, err := filepath.Rel(pn.homeDir, storagePath)
        if err == nil {
            return "~/" + rel, nil
        }
    }
    
    // If all else fails, return the absolute path
    return storagePath, nil
}

// GetCurrentDir returns the current working directory
func (pn *PathNormalizer) GetCurrentDir() string {
    return pn.currentDir
}

// GetHomeDir returns the user's home directory
func (pn *PathNormalizer) GetHomeDir() string {
    return pn.homeDir
}

// ValidateStoragePath checks if a stored path is valid
func (pn *PathNormalizer) ValidateStoragePath(path string) bool {
    // Must be absolute
    if !filepath.IsAbs(path) {
        return false
    }
    
    // Must be clean
    if filepath.Clean(path) != path {
        return false
    }
    
    return true
}
