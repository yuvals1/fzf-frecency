package finder

import (
    "fmt"
    "sort"

    "github.com/yuvals1/fzf-frecency/internal/frecency"
    "github.com/yuvals1/fzf-frecency/internal/log"
    "github.com/yuvals1/fzf-frecency/internal/pathutil"
)

// ScoredFile represents a file with its frecency score
type ScoredFile struct {
    Path  string
    Score int
}

// ScoredFiles is a slice of ScoredFile that can be sorted
type ScoredFiles []ScoredFile

func (sf ScoredFiles) Len() int      { return len(sf) }
func (sf ScoredFiles) Swap(i, j int) { sf[i], sf[j] = sf[j], sf[i] }
func (sf ScoredFiles) Less(i, j int) bool {
    // Sort by score (highest first)
    if sf[i].Score != sf[j].Score {
        return sf[i].Score > sf[j].Score
    }
    // Then by path (alphabetically)
    return sf[i].Path < sf[j].Path
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
    // Get all files first
    files, err := sf.FindFiles(root)
    if err != nil {
        return nil, fmt.Errorf("finding files: %w", err)
    }

    // Collect and score all files
    var scoredFiles ScoredFiles

    // Process files
    for path := range files {
        // First convert to storage path for scoring
        storagePath, err := sf.normalizer.ToStoragePath(path)
        if err != nil {
            log.Debug("Warning: couldn't convert path %s: %v", path, err)
            continue
        }

        // Get the score using the storage path
        score := sf.scorer.GetScore(storagePath)
        log.Debug("Got score %d for path: %s", score, storagePath)

        // Convert back to display path
        displayPath := path // default to original path
        if rel, err := sf.normalizer.ToDisplayPath(storagePath); err == nil {
            displayPath = rel
        }

        // Add to scored files with display path
        scoredFiles = append(scoredFiles, ScoredFile{
            Path:  displayPath,
            Score: score,
        })
    }

    // Sort the files
    sort.Sort(scoredFiles)

    // Create output channel
    out := make(chan ScoredFile)

    // Send sorted files through channel
    go func() {
        defer close(out)
        for _, file := range scoredFiles {
            out <- file
        }
    }()

    return out, nil
}
