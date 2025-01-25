package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/yuvals1/fzf-frecency/internal/finder"
    "github.com/yuvals1/fzf-frecency/internal/frecency"
    "github.com/yuvals1/fzf-frecency/internal/fzf"
    "github.com/yuvals1/fzf-frecency/internal/pathutil"
)

func main() {
    // Initialize path normalizer
    normalizer, err := pathutil.NewPathNormalizer()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating normalizer: %v\n", err)
        os.Exit(1)
    }

    // Create frecency scorer
    dataFile := filepath.Join(normalizer.GetHomeDir(), ".fzf_frecency.json")
    scorer, err := frecency.NewFrecencyScore(dataFile, normalizer)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating scorer: %v\n", err)
        os.Exit(1)
    }

    // Migrate records and prune missing files
    if err := scorer.MigrateRecords(); err != nil {
        fmt.Fprintf(os.Stderr, "Warning: Error during migration: %v\n", err)
    }
    if err := scorer.PruneMissingFiles(); err != nil {
        fmt.Fprintf(os.Stderr, "Warning: Error during pruning: %v\n", err)
    }

    // Create scored finder
    scoredFinder := finder.NewScoredFinder(scorer, normalizer)

    // Get current directory
    dir, err := os.Getwd()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
        os.Exit(1)
    }

    // Find scored files
    filesChan, err := scoredFinder.FindScoredFiles(dir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error finding files: %v\n", err)
        os.Exit(1)
    }

    // Setup FZF options
    opts := fzf.DefaultOptions()

    // Run FZF
    selected, err := fzf.RunFzf(filesChan, opts)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error running fzf: %v\n", err)
        os.Exit(1)
    }

    // Handle selected files
    for _, path := range selected {
        if err := scorer.UpdateAccess(path); err != nil {
            fmt.Fprintf(os.Stderr, "Error updating frecency: %v\n", err)
            continue
        }
        
        // Print selected path for shell script
        fmt.Println(path)
    }
}
