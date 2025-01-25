// main.go
package main

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
    // Clean the path first
    path = filepath.Clean(path)
    
    // Convert to absolute path if relative
    if !filepath.IsAbs(path) {
        path = filepath.Join(pn.currentDir, path)
    }
    
    // Try to make it relative to current directory
    relToCurrent, err := filepath.Rel(pn.currentDir, path)
    if err == nil && !strings.HasPrefix(relToCurrent, "..") {
        return relToCurrent, nil
    }
    
    // If path is under home directory, make it relative to home
    if strings.HasPrefix(path, pn.homeDir) {
        relToHome, err := filepath.Rel(pn.homeDir, path)
        if err == nil {
            return relToHome, nil
        }
    }
    
    // If we can't make it relative, return the cleaned absolute path
    return path, nil
}

func main() {
    normalizer, err := NewPathNormalizer()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating normalizer: %v\n", err)
        os.Exit(1)
    }

    // Test cases
    testPaths := []string{
        "test.txt",
        "./test.txt",
        "../test.txt",
        "/absolute/path/test.txt",
    }

    fmt.Println("Current directory:", normalizer.currentDir)
    fmt.Println("Home directory:", normalizer.homeDir)
    fmt.Println("\nPath normalization tests:")
    
    for _, path := range testPaths {
        normalized, err := normalizer.NormalizePath(path)
        if err != nil {
            fmt.Printf("Error normalizing %s: %v\n", path, err)
            continue
        }
        fmt.Printf("Original: %-20s -> Normalized: %s\n", path, normalized)
    }
}
