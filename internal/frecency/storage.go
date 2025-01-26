package frecency

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

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
    fs.records = make(map[string]*FileRecord)
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
