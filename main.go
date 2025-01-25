// main.go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// FileRecord represents a single file access record
type FileRecord struct {
    Path        string    `json:"path"`
    AccessCount int       `json:"access_count"`
    LastAccess  time.Time `json:"last_access"`
}

// FrecencyScore manages file access records and scoring
type FrecencyScore struct {
    records  map[string]*FileRecord
    dataFile string
}

// NewFrecencyScore creates a new FrecencyScore instance
func NewFrecencyScore(dataFile string) (*FrecencyScore, error) {
    fs := &FrecencyScore{
        records:  make(map[string]*FileRecord),
        dataFile: dataFile,
    }
    
    err := fs.Load()
    if err != nil && !os.IsNotExist(err) {
        return nil, fmt.Errorf("loading frecency data: %w", err)
    }
    
    return fs, nil
}

// Load reads frecency data from the data file
func (fs *FrecencyScore) Load() error {
    data, err := os.ReadFile(fs.dataFile)
    if err != nil {
        return err
    }
    
    var records []*FileRecord
    if err := json.Unmarshal(data, &records); err != nil {
        return fmt.Errorf("parsing frecency data: %w", err)
    }
    
    // Convert slice to map
    for _, record := range records {
        fs.records[record.Path] = record
    }
    
    return nil
}

// Save writes frecency data to the data file
func (fs *FrecencyScore) Save() error {
    // Convert map to slice for JSON serialization
    records := make([]*FileRecord, 0, len(fs.records))
    for _, record := range fs.records {
        records = append(records, record)
    }
    
    data, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return fmt.Errorf("serializing frecency data: %w", err)
    }
    
    // Create parent directory if it doesn't exist
    if err := os.MkdirAll(filepath.Dir(fs.dataFile), 0755); err != nil {
        return fmt.Errorf("creating data directory: %w", err)
    }
    
    // Write to temporary file first
    tmpFile := fs.dataFile + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return fmt.Errorf("writing temporary file: %w", err)
    }
    
    // Atomic rename
    if err := os.Rename(tmpFile, fs.dataFile); err != nil {
        os.Remove(tmpFile) // Clean up temp file
        return fmt.Errorf("renaming temporary file: %w", err)
    }
    
    return nil
}

func (fs *FrecencyScore) calculateScore(record *FileRecord) int {
    if record == nil {
        return 0
    }
    
    daysSince := int(time.Since(record.LastAccess).Hours()/24) + 1
    return (record.AccessCount * 100) / daysSince
}

func (fs *FrecencyScore) GetScore(path string) int {
    record, exists := fs.records[path]
    if !exists {
        return 0
    }
    return fs.calculateScore(record)
}

func (fs *FrecencyScore) UpdateAccess(path string) error {
    record, exists := fs.records[path]
    if !exists {
        record = &FileRecord{
            Path:        path,
            AccessCount: 0,
            LastAccess:  time.Now(),
        }
        fs.records[path] = record
    }
    
    record.AccessCount++
    record.LastAccess = time.Now()
    
    return fs.Save()
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

func main() {
    // Initialize path normalizer
    normalizer, err := NewPathNormalizer()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating normalizer: %v\n", err)
        os.Exit(1)
    }

    // Get home directory for data file location
    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
        os.Exit(1)
    }

    // Create frecency scorer with data file in home directory
    dataFile := filepath.Join(homeDir, ".fzf_frecency.json")
    scorer, err := NewFrecencyScore(dataFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating scorer: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Using frecency data file: %s\n\n", dataFile)

    // Test persistence
    testPaths := []string{
        "test.txt",
        "./docs/readme.md",
        "../parent/file.txt",
    }

    fmt.Println("1. Current scores from saved data:")
    for _, path := range testPaths {
        normalized, _ := normalizer.NormalizePath(path)
        score := scorer.GetScore(normalized)
        fmt.Printf("%-40s -> Score: %d\n", normalized, score)
    }

    fmt.Println("\n2. Updating access counts...")
    for _, path := range testPaths[:2] { // Access only first two files
        normalized, _ := normalizer.NormalizePath(path)
        if err := scorer.UpdateAccess(normalized); err != nil {
            fmt.Fprintf(os.Stderr, "Error updating access: %v\n", err)
            continue
        }
        fmt.Printf("Accessed: %s\n", normalized)
    }

    fmt.Println("\n3. Final scores (should persist after restart):")
    for _, path := range testPaths {
        normalized, _ := normalizer.NormalizePath(path)
        score := scorer.GetScore(normalized)
        fmt.Printf("%-40s -> Score: %d\n", normalized, score)
    }
}
