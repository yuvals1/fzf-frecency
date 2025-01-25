package frecency

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

// GetScore returns the frecency score for a path
func (fs *FrecencyScore) GetScore(path string) int {
    record, exists := fs.records[path]
    if !exists {
        return 0
    }
    return record.CalculateScore()
}

// UpdateAccess records a new access for a path
func (fs *FrecencyScore) UpdateAccess(path string) error {
    record, exists := fs.records[path]
    if !exists {
        record = NewFileRecord(path)
        fs.records[path] = record
    }
    
    record.UpdateAccess()
    return fs.Save()
}

// GetAllScores returns all paths and their current scores
func (fs *FrecencyScore) GetAllScores() map[string]int {
    scores := make(map[string]int)
    for path, record := range fs.records {
        scores[path] = record.CalculateScore()
    }
    return scores
}
