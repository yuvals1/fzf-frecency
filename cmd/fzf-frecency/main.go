package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "bufio"
    "strings"

    "github.com/yuvals1/fzf-frecency/internal/finder"
    "github.com/yuvals1/fzf-frecency/internal/frecency"
    "github.com/yuvals1/fzf-frecency/internal/fzf"
    "github.com/yuvals1/fzf-frecency/internal/pathutil"
    "github.com/yuvals1/fzf-frecency/internal/log"
)

type config struct {
    usePresetDirs bool
}

func readPresetPaths() ([]string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("getting home directory: %w", err)
    }

    presetFile := filepath.Join(homeDir, ".fzf_preset_paths")
    file, err := os.Open(presetFile)
    if err != nil {
        return nil, fmt.Errorf("opening preset paths file %s: %w", presetFile, err)
    }
    defer file.Close()

    var paths []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        path := strings.TrimSpace(scanner.Text())
        if path != "" {
            paths = append(paths, path)
            log.Debug("Read preset path: %s", path)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("reading preset paths: %w", err)
    }

    return paths, nil
}

func main() {
    // Parse flags
    cfg := config{}
    flag.BoolVar(&cfg.usePresetDirs, "use-preset-dirs", false, "Use directories from preset paths file")
    flag.Parse()

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

    // Create base finder
    fileFinder := finder.NewFileFinder()
    
    // Set exclude patterns to match user's configuration
    fileFinder.SetExcludePatterns([]string{
        "*.mypy",
        "*.mypy_cache",
        "*.git",
    })

    // Create scored finder
    scoredFinder := finder.NewScoredFinder(scorer, normalizer)

    // Get current directory
    dir, err := os.Getwd()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
        os.Exit(1)
    }

    // Find scored files based on mode
    var filesChan <-chan finder.ScoredFile
    if cfg.usePresetDirs {
        // Get preset paths
        presetPaths, err := readPresetPaths()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading preset paths: %v\n", err)
            os.Exit(1)
        }
        
        // Add current directory to the list
        searchPaths := append([]string{dir}, presetPaths...)
        log.Debug("Searching in directories:")
        for _, path := range searchPaths {
            log.Debug("  %s", path)
        }

        // Use multi-directory search
        filesChan, err = scoredFinder.FindScoredFilesInDirs(searchPaths)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error finding files in multiple directories: %v\n", err)
            os.Exit(1)
        }
    } else {
        // Use single directory search (original behavior)
        filesChan, err = scoredFinder.FindScoredFiles(dir)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error finding files: %v\n", err)
            os.Exit(1)
        }
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
