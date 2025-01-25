// main.go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// FileRecord represents a single file access record
type FileRecord struct {
    Path       string    // Normalized path
    AccessCount int      // Number of times accessed
    LastAccess time.Time // Last access timestamp
}

// FrecencyScore manages file access records and scoring
type FrecencyScore struct {
    records map[string]*FileRecord
}

// NewFrecencyScore creates a new FrecencyScore instance
func NewFrecencyScore() *FrecencyScore {
    return &FrecencyScore{
        records: make(map[string]*FileRecord),
    }
}

// calculateScore computes the frecency score for a file
// Score = (access_count * 100) / days_since_last_access
func (fs *FrecencyScore) calculateScore(record *FileRecord) int {
    if record == nil {
        return 0
    }
    
    daysSince := int(time.Since(record.LastAccess).Hours()/24) + 1
    return (record.AccessCount * 100) / daysSince
}

// GetScore returns the frecency score for a path
func (fs *FrecencyScore) GetScore(path string) int {
    record, exists := fs.records[path]
    if !exists {
        return 0
    }
    return fs.calculateScore(record)
}

// UpdateAccess records a new access for a path
func (fs *FrecencyScore) UpdateAccess(path string) {
    record, exists := fs.records[path]
    if !exists {
        record = &FileRecord{
            Path:       path,
            AccessCount: 0,
            LastAccess: time.Now(),
        }
        fs.records[path] = record
    }
    
    record.AccessCount++
    record.LastAccess = time.Now()
}

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
    
    // If path is already relative and contains .., return it as is
    if !filepath.IsAbs(path) && strings.HasPrefix(path, "..") {
        return path, nil
    }
    
    // Convert to absolute path if relative
    absPath := path
    if !filepath.IsAbs(path) {
        absPath = filepath.Join(pn.currentDir, path)
    }
    
    // Try to make it relative to current directory
    relToCurrent, err := filepath.Rel(pn.currentDir, absPath)
    if err == nil && !strings.HasPrefix(relToCurrent, "..") {
        return relToCurrent, nil
    }
    
    // If path is under home directory, make it relative to home
    if strings.HasPrefix(absPath, pn.homeDir) {
        relToHome, err := filepath.Rel(pn.homeDir, absPath)
        if err == nil {
            return relToHome, nil
        }
    }
    
    // If we can't make it relative, return the cleaned absolute path
    return absPath, nil
}

func main() {
    normalizer, err := NewPathNormalizer()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating normalizer: %v\n", err)
        os.Exit(1)
    }

    // Create a new frecency scorer
    scorer := NewFrecencyScore()

    // Test scoring
    testPaths := []string{
        "test.txt",
        "./docs/readme.md",
        "../parent/file.txt",
        "/absolute/path/test.txt",
        filepath.Join(normalizer.homeDir, "Documents/test.txt"),
    }

    fmt.Println("Frecency scoring tests:")
    fmt.Println("1. Initial scores (should all be 0):")
    for _, path := range testPaths {
        normalized, _ := normalizer.NormalizePath(path)
        score := scorer.GetScore(normalized)
        fmt.Printf("%-40s -> Score: %d\n", normalized, score)
    }

    fmt.Println("\n2. After accessing some files:")
    // Simulate some file accesses
    accessPaths := []string{
        "test.txt",                    // Access once
        "./docs/readme.md",            // Access twice
        "./docs/readme.md",
    }

    for _, path := range accessPaths {
        normalized, _ := normalizer.NormalizePath(path)
        scorer.UpdateAccess(normalized)
    }

    // Print updated scores
    for _, path := range testPaths {
        normalized, _ := normalizer.NormalizePath(path)
        score := scorer.GetScore(normalized)
        fmt.Printf("%-40s -> Score: %d\n", normalized, score)
    }
}
