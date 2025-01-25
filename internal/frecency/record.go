package frecency

import (
    "time"
)

// FileRecord represents a single file access record
type FileRecord struct {
    Path        string    `json:"path"`
    AccessCount int       `json:"access_count"`
    LastAccess  time.Time `json:"last_access"`
}

// NewFileRecord creates a new FileRecord for the given path
func NewFileRecord(path string) *FileRecord {
    return &FileRecord{
        Path:        path,
        AccessCount: 0,
        LastAccess:  time.Now(),
    }
}

// UpdateAccess increments the access count and updates the last access time
func (r *FileRecord) UpdateAccess() {
    r.AccessCount++
    r.LastAccess = time.Now()
}

// CalculateScore computes the frecency score for this record
// Score = (access_count * 100) / days_since_last_access
func (r *FileRecord) CalculateScore() int {
    daysSince := int(time.Since(r.LastAccess).Hours()/24) + 1
    return (r.AccessCount * 100) / daysSince
}
