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

// NormalizePath converts a path to be relative to the current directory
func (pn *PathNormalizer) NormalizePath(path string) (string, error) {
    path = filepath.Clean(path)
    
    if !filepath.IsAbs(path) && strings.HasPrefix(path, "..") {
        return path, nil
    }
    
    absPath := path
    if !filepath.IsAbs(path) {
        absPath = filepath.Join(pn.currentDir, path)
    }
    
    relToCurrent, err := filepath.Rel(pn.currentDir, absPath)
    if err == nil && !strings.HasPrefix(relToCurrent, "..") {
        return relToCurrent, nil
    }
    
    if strings.HasPrefix(absPath, pn.homeDir) {
        relToHome, err := filepath.Rel(pn.homeDir, absPath)
        if err == nil {
            return relToHome, nil
        }
    }
    
    return absPath, nil
}

// GetCurrentDir returns the current working directory
func (pn *PathNormalizer) GetCurrentDir() string {
    return pn.currentDir
}

// GetHomeDir returns the user's home directory
func (pn *PathNormalizer) GetHomeDir() string {
    return pn.homeDir
}
