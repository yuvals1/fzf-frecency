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
    scorer, err := frecency.NewFrecencyScore(dataFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating scorer: %v\n", err)
        os.Exit(1)
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

    // Setup FZF options with default preview
    opts := fzf.DefaultOptions()

    // Run FZF
    selected, err := fzf.RunFzf(filesChan, opts)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error running fzf: %v\n", err)
        os.Exit(1)
    }

    // Handle selected files
    for _, path := range selected {
        // Update frecency score for selected file
        normalized, err := normalizer.NormalizePath(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error normalizing path: %v\n", err)
            continue
        }
        
        if err := scorer.UpdateAccess(normalized); err != nil {
            fmt.Fprintf(os.Stderr, "Error updating frecency: %v\n", err)
            continue
        }
        
        // Print selected path (can be used by shell script)
        fmt.Println(path)
    }
}
