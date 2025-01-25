package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/yuvals1/fzf-frecency/internal/frecency"
    "github.com/yuvals1/fzf-frecency/internal/pathutil"
)

func main() {
    // Initialize path normalizer
    normalizer, err := pathutil.NewPathNormalizer()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating normalizer: %v\n", err)
        os.Exit(1)
    }

    // Create frecency scorer with data file in home directory
    dataFile := filepath.Join(normalizer.GetHomeDir(), ".fzf_frecency.json")
    scorer, err := frecency.NewFrecencyScore(dataFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating scorer: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Using frecency data file: %s\n\n", dataFile)

    // Test paths
    testPaths := []string{
        "test.txt",
        "./docs/readme.md",
        "../parent/file.txt",
    }

    fmt.Println("1. Current scores from saved data:")
    for _, path := range testPaths {
        normalized, err := normalizer.NormalizePath(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error normalizing path: %v\n", err)
            continue
        }
        score := scorer.GetScore(normalized)
        fmt.Printf("%-40s -> Score: %d\n", normalized, score)
    }

    fmt.Println("\n2. Updating access counts...")
    for _, path := range testPaths[:2] {
        normalized, err := normalizer.NormalizePath(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error normalizing path: %v\n", err)
            continue
        }
        if err := scorer.UpdateAccess(normalized); err != nil {
            fmt.Fprintf(os.Stderr, "Error updating access: %v\n", err)
            continue
        }
        fmt.Printf("Accessed: %s\n", normalized)
    }

    fmt.Println("\n3. Final scores (should persist after restart):")
    for _, path := range testPaths {
        normalized, err := normalizer.NormalizePath(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error normalizing path: %v\n", err)
            continue
        }
        score := scorer.GetScore(normalized)
        fmt.Printf("%-40s -> Score: %d\n", normalized, score)
    }
}
