package fzf

import (
    "fmt"
    "io"
    "os/exec"
    "strings"

    "github.com/yuvals1/fzf-frecency/internal/finder"
)

// FzfOptions contains configuration for fzf
type FzfOptions struct {
    Preview     string   // Preview command
    Height      string   // Window height (e.g., "50%")
    MultiSelect bool     // Allow multiple selection
    ExtraArgs   []string // Additional fzf arguments
}

// DefaultOptions returns default FZF options
func DefaultOptions() *FzfOptions {
    return &FzfOptions{
        Preview:     "bat --color=always {2..}", // Use {2..} to get everything after the tab
        Height:      "50%",
        MultiSelect: false,
        ExtraArgs: []string{
            "--ansi",
            "--border",
            "--reverse",
            "--delimiter=\t",
            "--with-nth=1,2", // Show both score and path
        },
    }
}

// FormatScoredFile formats a scored file for FZF display
func FormatScoredFile(file finder.ScoredFile) string {
    return fmt.Sprintf("%5d\t%s", file.Score, file.Path)
}

// RunFzf runs fzf with the provided scored files
func RunFzf(files <-chan finder.ScoredFile, opts *FzfOptions) ([]string, error) {
    if opts == nil {
        opts = DefaultOptions()
    }

    // Build fzf command
    args := []string{
        "--height", opts.Height,
        "--preview", opts.Preview,
    }
    
    if opts.MultiSelect {
        args = append(args, "-m")
    }
    
    args = append(args, opts.ExtraArgs...)
    
    // Start fzf process
    cmd := exec.Command("fzf", args...)
    
    // Create pipes for stdin and stdout
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, fmt.Errorf("creating stdin pipe: %w", err)
    }
    
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("creating stdout pipe: %w", err)
    }
    
    // Start fzf before writing to stdin
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("starting fzf: %w", err)
    }

    // Write scored files to fzf
    go func() {
        defer stdin.Close()
        for file := range files {
            fmt.Fprintln(stdin, FormatScoredFile(file))
        }
    }()

    // Read fzf output
    output, err := io.ReadAll(stdout)
    if err != nil {
        return nil, fmt.Errorf("reading fzf output: %w", err)
    }

    // Wait for fzf to finish
    if err := cmd.Wait(); err != nil {
        // Exit code 130 means user interrupted (e.g., by pressing Esc)
        if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
            return nil, nil
        }
        return nil, fmt.Errorf("fzf process: %w", err)
    }

    // Process output (split by newlines and extract paths)
    var results []string
    for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
        if line == "" {
            continue
        }
        // Split on tab and take the path part (after the score)
        parts := strings.SplitN(line, "\t", 2)
        if len(parts) == 2 {
            results = append(results, parts[1])
        }
    }

    return results, nil
}
