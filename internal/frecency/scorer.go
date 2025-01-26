package frecency

import (
    "fmt"
    "os"

    "github.com/yuvals1/fzf-frecency/internal/pathutil"
)

// FrecencyScore manages file access records and scoring
type FrecencyScore struct {
    records    map[string]*FileRecord
    dataFile   string
    normalizer *pathutil.PathNormalizer
}

// NewFrecencyScore creates a new FrecencyScore instance
func NewFrecencyScore(dataFile string, normalizer *pathutil.PathNormalizer) (*FrecencyScore, error) {
    fs := &FrecencyScore{
        records:    make(map[string]*FileRecord),
        dataFile:   dataFile,
        normalizer: normalizer,
    }
    
    err := fs.Load()
    if err != nil && !os.IsNotExist(err) {
        return nil, fmt.Errorf("loading frecency data: %w", err)
    }
    
    return fs, nil
}

// GetScore returns the frecency score for a path
func (fs *FrecencyScore) GetScore(path string) int {
    if record, exists := fs.records[path]; exists {
        return record.CalculateScore()
    }
    return 0
}

// UpdateAccess records a new access for a path
func (fs *FrecencyScore) UpdateAccess(path string) error {
    storagePath, err := fs.normalizer.ToStoragePath(path)
    if err != nil {
        return fmt.Errorf("converting to storage path: %w", err)
    }
    
    record, exists := fs.records[storagePath]
    if !exists {
        record = NewFileRecord(storagePath)
        fs.records[storagePath] = record
    }
    
    record.UpdateAccess()
    return fs.Save()
}

// MigrateRecords updates all stored paths to absolute format
func (fs *FrecencyScore) MigrateRecords() error {
    newRecords := make(map[string]*FileRecord)
    
    for oldPath, record := range fs.records {
        // Convert to absolute path
        newPath, err := fs.normalizer.ToStoragePath(oldPath)
        if err != nil {
            continue // Skip invalid paths
        }
        
        // If we already have this path, merge the records
        if existingRecord, exists := newRecords[newPath]; exists {
            if record.AccessCount > existingRecord.AccessCount {
                existingRecord.AccessCount = record.AccessCount
            }
            if record.LastAccess.After(existingRecord.LastAccess) {
                existingRecord.LastAccess = record.LastAccess
            }
        } else {
            record.Path = newPath
            newRecords[newPath] = record
        }
    }
    
    fs.records = newRecords
    return fs.Save()
}
