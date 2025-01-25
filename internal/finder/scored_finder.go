package finder

import (
    "fmt"

    "github.com/yuvals1/fzf-frecency/internal/frecency"
    "github.com/yuvals1/fzf-frecency/internal/pathutil"
)

// ScoredFile represents a file with its frecency score
type ScoredFile struct {
    Path  string
    Score int
}

// ScoredFinder combines file finding with frecency scoring
type ScoredFinder struct {
    *FileFinder
    scorer     *frecency.FrecencyScore
    normalizer *pathutil.PathNormalizer
}

// NewScoredFinder creates a new ScoredFinder
func NewScoredFinder(scorer *frecency.FrecencyScore, normalizer *pathutil.PathNormalizer) *ScoredFinder {
    return &ScoredFinder{
        FileFinder: NewFileFinder(),
        scorer:     scorer,
        normalizer: normalizer,
    }
}

// FindScoredFiles returns a channel of files with their frecency scores
func (sf *ScoredFinder) FindScoredFiles(root string) (<-chan ScoredFile, error) {
    // Create output channel for scored files
    scoredFiles := make(chan ScoredFile)
    
    // Get raw files channel
    files, err := sf.FindFiles(root)
    if err != nil {
        return nil, fmt.Errorf("finding files: %w", err)
    }

    // Process files and add scores
    go func() {
        defer close(scoredFiles)
        
        for path := range files {
            // Normalize the path for scoring
            normalizedPath, err := sf.normalizer.NormalizePath(path)
            if err != nil {
                fmt.Printf("Warning: couldn't normalize path %s: %v\n", path, err)
                continue
            }

            // Get score and send scored file
            score := sf.scorer.GetScore(normalizedPath)
            scoredFiles <- ScoredFile{
                Path:  path,
                Score: score,
            }
        }
    }()

    return scoredFiles, nil
}
