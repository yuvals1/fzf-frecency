package finder

import (
    "fmt"
    "sort"

    "github.com/yuvals1/fzf-frecency/internal/frecency"
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
        // Normalize the path for scoring
        normalizedPath, err := sf.normalizer.NormalizePath(path)
        if err != nil {
            fmt.Printf("Warning: couldn't normalize path %s: %v\n", path, err)
            continue
        }

        // Get score and add to slice
        score := sf.scorer.GetScore(normalizedPath)
        scoredFiles = append(scoredFiles, ScoredFile{
            Path:  path,
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
