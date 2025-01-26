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
    RawPath  string  // Original unmodified path for operations
    Path     string  // Path for display
    Score    int
}

type ScoredFiles []ScoredFile

func (sf ScoredFiles) Len() int      { return len(sf) }
func (sf ScoredFiles) Swap(i, j int) { sf[i], sf[j] = sf[j], sf[i] }
func (sf ScoredFiles) Less(i, j int) bool {
    if sf[i].Score != sf[j].Score {
        return sf[i].Score > sf[j].Score
    }
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

// processFile converts a file path to a ScoredFile
func (sf *ScoredFinder) processFile(path string) (ScoredFile, error) {
    // Convert to storage path for scoring
    storagePath, err := sf.normalizer.ToStoragePath(path)
    if err != nil {
        return ScoredFile{}, fmt.Errorf("converting to storage path: %w", err)
    }

    // Get the score using the storage path
    score := sf.scorer.GetScore(storagePath)
    log.Debug("Got score %d for storage path: %s", score, storagePath)

    // Keep both raw and display paths
    displayPath := path
    if rel, err := sf.normalizer.ToDisplayPath(storagePath); err == nil {
        if !sf.normalizer.IsOutsideCurrentDir(rel) {
            displayPath = rel
        }
    }

    return ScoredFile{
        RawPath:  path,        // Original path for operations
        Path:     displayPath, // Path for display
        Score:    score,
    }, nil
}

// FindScoredFiles returns a channel of files with their frecency scores from a single directory
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
        scoredFile, err := sf.processFile(path)
        if err != nil {
            log.Debug("Warning: couldn't process file %s: %v", path, err)
            continue
        }
        scoredFiles = append(scoredFiles, scoredFile)
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

// FindScoredFilesInDirs returns a channel of scored files from multiple directories
func (sf *ScoredFinder) FindScoredFilesInDirs(roots []string) (<-chan ScoredFile, error) {
    // Get all files from all directories
    filesChan, err := sf.FileFinder.FindFilesInDirs(roots)
    if err != nil {
        return nil, fmt.Errorf("finding files in directories: %w", err)
    }

    // Use a map to deduplicate files by storage path
    seen := make(map[string]bool)
    var scoredFiles ScoredFiles

    // Process all files
    for path := range filesChan {
        // Get storage path for deduplication
        storagePath, err := sf.normalizer.ToStoragePath(path)
        if err != nil {
            log.Debug("Warning: couldn't get storage path for %s: %v", path, err)
            continue
        }

        // Skip if we've already seen this file
        if seen[storagePath] {
            log.Debug("Skipping duplicate file: %s", storagePath)
            continue
        }
        seen[storagePath] = true

        // Process the file
        scoredFile, err := sf.processFile(path)
        if err != nil {
            log.Debug("Warning: couldn't process file %s: %v", path, err)
            continue
        }
        scoredFiles = append(scoredFiles, scoredFile)
    }

    // Sort all files
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
